package action

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/promhippie/prometheus-vcd-sd/pkg/client"
)

var (
	// providerPrefix defines the general prefix for all labels.
	providerPrefix = model.MetaLabelPrefix + "vcd_"

	// Labels defines all available labels for this provider.
	Labels = map[string]string{
		"id":                providerPrefix + "id",
		"metadataPrefix":    providerPrefix + "metadata_",
		"name":              providerPrefix + "name",
		"networkPrefix":     providerPrefix + "network_",
		"numCoresPerSocket": providerPrefix + "num_cores_per_socket",
		"numCpus":           providerPrefix + "num_cpus",
		"org":               providerPrefix + "org",
		"osType":            providerPrefix + "os_type",
		"project":           providerPrefix + "project",
		"status":            providerPrefix + "status",
		"storageProfile":    providerPrefix + "storage_profile",
		"vdc":               providerPrefix + "vdc",
	}

	// ErrClientEndpoint defines an error if the client auth fails.
	ErrClientEndpoint = errors.New("failed to parse api url")
)

// Discoverer implements the Prometheus discoverer interface.
type Discoverer struct {
	configs   map[string]*client.Client
	logger    log.Logger
	refresh   int
	separator string
	lasts     map[string]struct{}
}

// Run initializes fetching the targets for service discovery.
func (d Discoverer) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	ticker := time.NewTicker(time.Duration(d.refresh) * time.Second)

	for {
		targets, err := d.getTargets(ctx)

		if err == nil {
			ch <- targets
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (d *Discoverer) getTargets(ctx context.Context) ([]*targetgroup.Group, error) {
	current := make(map[string]struct{})
	targets := make([]*targetgroup.Group, 0)

	for project, config := range d.configs {
		if err := config.Authenticate(); err != nil {
			level.Warn(d.logger).Log(
				"msg", "Failed to authenticate",
				"project", project,
				"err", err,
			)

			requestFailures.WithLabelValues(project, "auth").Inc()
			continue
		}

		defer config.Disconnect()

		nowOrg := time.Now()
		org, err := config.Upstream.GetOrgByNameOrId(config.Organization)
		requestDuration.WithLabelValues(project, "org").Observe(time.Since(nowOrg).Seconds())

		if err != nil {
			level.Warn(d.logger).Log(
				"msg", "Failed to fetch org",
				"project", project,
				"err", err,
			)

			requestFailures.WithLabelValues(project, "org").Inc()
			continue
		}

		nowVdc := time.Now()
		vdc, err := org.GetVDCByNameOrId(config.Datacenter, false)
		requestDuration.WithLabelValues(project, "vdc").Observe(time.Since(nowVdc).Seconds())

		if err != nil {
			level.Warn(d.logger).Log(
				"msg", "Failed to fetch vdc",
				"project", project,
				"err", err,
			)

			requestFailures.WithLabelValues(project, "vdc").Inc()
			continue
		}

		vappNames := []string{}

		for _, entities := range vdc.Vdc.ResourceEntities {
			for _, entity := range entities.ResourceEntity {
				if entity.Type == "application/vnd.vmware.vcloud.vApp+xml" {
					vappNames = append(vappNames, entity.Name)
				}
			}
		}

		for _, vappName := range vappNames {
			nowVapp := time.Now()
			vapp, err := vdc.GetVAppByName(vappName, false)
			requestDuration.WithLabelValues(project, "vapp").Observe(time.Since(nowVapp).Seconds())

			if err != nil {
				level.Warn(d.logger).Log(
					"msg", "Failed to fetch servers",
					"project", project,
					"vapp", vappName,
					"err", err,
				)

				requestFailures.WithLabelValues(project, "vapp").Inc()
				continue
			}

			if vapp.VApp.Children == nil {
				level.Debug(d.logger).Log(
					"msg", "No servers defined",
					"project", project,
					"vapp", vappName,
				)

				continue
			}

			servers := vapp.VApp.Children.VM

			level.Debug(d.logger).Log(
				"msg", "Requested servers",
				"project", project,
				"vapp", vappName,
				"count", len(servers),
			)

			for _, server := range servers {
				if len(server.NetworkConnectionSection.NetworkConnection) < 1 {
					continue
				}

				nowMachine := time.Now()
				vm, err := vapp.GetVMByName(server.Name, false)
				requestDuration.WithLabelValues(project, "vm").Observe(time.Since(nowMachine).Seconds())

				if err != nil {
					level.Warn(d.logger).Log(
						"msg", "Failed to fetch vm",
						"project", project,
						"vapp", vappName,
						"server", server.Name,
						"id", server.ID,
						"err", err,
					)

					requestFailures.WithLabelValues(project, "vm").Inc()
					continue
				}

				nowMeta := time.Now()
				metadata, err := vm.GetMetadata()
				requestDuration.WithLabelValues(project, "metadata").Observe(time.Since(nowMeta).Seconds())

				if err != nil {
					level.Warn(d.logger).Log(
						"msg", "Failed to fetch metadata",
						"project", project,
						"vapp", vappName,
						"server", server.Name,
						"id", server.ID,
						"err", err,
					)

					requestFailures.WithLabelValues(project, "metadata").Inc()
					continue
				}

				target := &targetgroup.Group{
					Source: fmt.Sprintf("vcd/%s", vm.VM.ID),
					Targets: []model.LabelSet{
						{
							model.AddressLabel: model.LabelValue(vm.VM.NetworkConnectionSection.NetworkConnection[0].IPAddress),
						},
					},
					Labels: model.LabelSet{
						model.AddressLabel:                 model.LabelValue(vm.VM.NetworkConnectionSection.NetworkConnection[0].IPAddress),
						model.LabelName(Labels["project"]): model.LabelValue(project),
						model.LabelName(Labels["org"]):     model.LabelValue(config.Organization),
						model.LabelName(Labels["vdc"]):     model.LabelValue(config.Datacenter),
						model.LabelName(Labels["name"]):    model.LabelValue(vm.VM.Name),
						model.LabelName(Labels["id"]):      model.LabelValue(vm.VM.ID),
						model.LabelName(Labels["status"]):  model.LabelValue(strconv.Itoa(vm.VM.Status)),
					},
				}

				if vm.VM.VmSpecSection != nil {
					target.Labels[model.LabelName(Labels["osType"])] = model.LabelValue(vm.VM.VmSpecSection.OsType)
				}

				if vm.VM.VmSpecSection != nil {
					target.Labels[model.LabelName(Labels["numCpus"])] = model.LabelValue(strconv.Itoa(*vm.VM.VmSpecSection.NumCpus))
					target.Labels[model.LabelName(Labels["numCoresPerSocket"])] = model.LabelValue(strconv.Itoa(*vm.VM.VmSpecSection.NumCoresPerSocket))
				}

				if vm.VM.StorageProfile != nil {
					target.Labels[model.LabelName(Labels["storageProfile"])] = model.LabelValue(vm.VM.StorageProfile.Name)
				}

				for _, network := range vm.VM.NetworkConnectionSection.NetworkConnection {
					normalized := normalizeLabel(network.Network)
					target.Labels[model.LabelName(Labels["networkPrefix"]+normalized)] = model.LabelValue(network.IPAddress)
				}

				for _, entry := range metadata.MetadataEntry {
					normalized := normalizeLabel(entry.Key)
					target.Labels[model.LabelName(Labels["metadataPrefix"]+normalized)] = model.LabelValue(entry.TypedValue.Value)
				}

				level.Debug(d.logger).Log(
					"msg", "Server added",
					"project", project,
					"vapp", vappName,
					"server", vm.VM.Name,
					"source", target.Source,
				)

				current[target.Source] = struct{}{}
				targets = append(targets, target)
			}
		}

		config.Disconnect()
	}

	for k := range d.lasts {
		if _, ok := current[k]; !ok {
			level.Debug(d.logger).Log(
				"msg", "Server deleted",
				"source", k,
			)

			targets = append(
				targets,
				&targetgroup.Group{
					Source: k,
				},
			)
		}
	}

	d.lasts = current
	return targets, nil
}

func normalizeLabel(val string) string {
	replaces := map[string]string{
		"-": "_",
		".": "_",
		",": "_",
	}

	for original, replaced := range replaces {
		val = strings.ReplaceAll(
			val,
			original,
			replaced,
		)
	}

	return strings.ToLower(
		val,
	)
}
