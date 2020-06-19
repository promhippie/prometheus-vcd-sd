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
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

const (
	vcdPrefix              = model.MetaLabelPrefix + "vcd_"
	projectLabel           = vcdPrefix + "project"
	orgLabel               = vcdPrefix + "org"
	vdcLabel               = vcdPrefix + "vdc"
	nameLabel              = vcdPrefix + "name"
	statusLabel            = vcdPrefix + "status"
	osTypeLabel            = vcdPrefix + "os_type"
	numCpusLabel           = vcdPrefix + "num_cpus"
	numCoresPerSocketLabel = vcdPrefix + "num_cores_per_socket"
	storageProfileLabel    = vcdPrefix + "storage_profile"
	networkPrefix          = vcdPrefix + "network_"
	metadataPrefix         = vcdPrefix + "metadata_"
)

var (
	// ErrClientEndpoint defines an error if the client auth fails.
	ErrClientEndpoint = errors.New("failed to parse api url")

	// ErrClientAuth defines an error if the client auth fails.
	ErrClientAuth = errors.New("failed to authenticate client")
)

// Config wraps the vCloud Director client including org and vdc names.
type Config struct {
	client *govcd.VCDClient
	org    string
	vdc    string
}

// Discoverer implements the Prometheus discoverer interface.
type Discoverer struct {
	configs   map[string]*Config
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
		// 	now := time.Now()
		org, err := config.client.GetOrgByName(config.org)
		// 	requestDuration.WithLabelValues(project).Observe(time.Since(now).Seconds())

		if err != nil {
			level.Warn(d.logger).Log(
				"msg", "Failed to fetch org",
				"project", project,
				"err", err,
			)

			requestFailures.WithLabelValues(project).Inc()
			continue
		}

		// 	now := time.Now()
		vdc, err := org.GetVDCByName(config.vdc, false)
		// 	requestDuration.WithLabelValues(project).Observe(time.Since(now).Seconds())

		if err != nil {
			level.Warn(d.logger).Log(
				"msg", "Failed to fetch vdc",
				"project", project,
				"err", err,
			)

			requestFailures.WithLabelValues(project).Inc()
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
			// 	now := time.Now()
			vapp, err := vdc.GetVAppByNameOrId(vappName, false)
			// 	requestDuration.WithLabelValues(project).Observe(time.Since(now).Seconds())

			if err != nil {
				level.Warn(d.logger).Log(
					"msg", "Failed to fetch servers",
					"project", project,
					"vapp", vappName,
					"err", err,
				)

				requestFailures.WithLabelValues(project).Inc()
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

				// // 	now := time.Now()
				// metadata, err := vapp.GetMetadata()
				// // 	requestDuration.WithLabelValues(project).Observe(time.Since(now).Seconds())

				// if err != nil {
				// 	level.Warn(d.logger).Log(
				// 		"msg", "Failed to fetch metadata",
				// 		"project", project,
				// 		"vapp", vappName,
				// 		"err", err,
				// 	)

				// 	requestFailures.WithLabelValues(project).Inc()
				// 	continue
				// }

				target := &targetgroup.Group{
					Source: fmt.Sprintf("vcd/%s", server.ID),
					Targets: []model.LabelSet{
						{
							model.AddressLabel: model.LabelValue(server.NetworkConnectionSection.NetworkConnection[0].IPAddress),
						},
					},
					Labels: model.LabelSet{
						model.AddressLabel:                      model.LabelValue(server.NetworkConnectionSection.NetworkConnection[0].IPAddress),
						model.LabelName(projectLabel):           model.LabelValue(project),
						model.LabelName(orgLabel):               model.LabelValue(config.org),
						model.LabelName(vdcLabel):               model.LabelValue(config.vdc),
						model.LabelName(nameLabel):              model.LabelValue(server.Name),
						model.LabelName(statusLabel):            model.LabelValue(strconv.Itoa(server.Status)),
						model.LabelName(osTypeLabel):            model.LabelValue(server.VmSpecSection.OsType),
						model.LabelName(numCpusLabel):           model.LabelValue(strconv.Itoa(*server.VmSpecSection.NumCpus)),
						model.LabelName(numCoresPerSocketLabel): model.LabelValue(strconv.Itoa(*server.VmSpecSection.NumCoresPerSocket)),
						model.LabelName(storageProfileLabel):    model.LabelValue(server.StorageProfile.Name),
					},
				}

				for _, network := range server.NetworkConnectionSection.NetworkConnection {
					target.Labels[model.LabelName(networkPrefix+strings.ToLower(network.Network))] = model.LabelValue(network.IPAddress)
				}

				// for _, entry := range metadata.MetadataEntry {
				// 	target.Labels[model.LabelName(metadataPrefix+strings.ToLower(entry.Key))] = model.LabelValue(entry.TypedValue.Value)
				// }

				level.Debug(d.logger).Log(
					"msg", "Server added",
					"project", project,
					"vapp", vappName,
					"server", server.Name,
					"source", target.Source,
				)

				current[target.Source] = struct{}{}
				targets = append(targets, target)
			}
		}
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
