package action

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

const (
	vcdPrefix    = model.MetaLabelPrefix + "vcd_"
	projectLabel = vcdPrefix + "project"
	nameLabel    = vcdPrefix + "name"
)

var (
	// ErrClientEndpoint defines an error if the client auth fails.
	ErrClientEndpoint = errors.New("failed to parse api url")

	// ErrClientAuth defines an error if the client auth fails.
	ErrClientAuth = errors.New("failed to authenticate client")
)

// Discoverer implements the Prometheus discoverer interface.
type Discoverer struct {
	clients   map[string]*govcd.VCDClient
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

	// for project, client := range d.clients {

	// 	now := time.Now()
	// 	servers, err := client.Server.All(ctx)
	// 	requestDuration.WithLabelValues(project).Observe(time.Since(now).Seconds())

	// 	if err != nil {
	// 		level.Warn(d.logger).Log(
	// 			"msg", "Failed to fetch servers",
	// 			"project", project,
	// 			"err", err,
	// 		)

	// 		requestFailures.WithLabelValues(project).Inc()
	// 		continue
	// 	}

	// 	level.Debug(d.logger).Log(
	// 		"msg", "Requested servers",
	// 		"project", project,
	// 		"count", len(servers),
	// 	)

	// 	for _, server := range servers {
	// 		var (
	// 			imageType string
	// 			imageName string
	// 			osFlavor  string
	// 			osVersion string
	// 		)

	// 		if server.Image != nil {
	// 			imageType = string(server.Image.Type)
	// 			imageName = server.Image.Name
	// 			osFlavor = server.Image.OSFlavor
	// 			osVersion = server.Image.OSVersion
	// 		}

	// 		target := &targetgroup.Group{
	// 			Source: fmt.Sprintf("vcd/%d", server.ID),
	// 			Targets: []model.LabelSet{
	// 				{
	// 					model.AddressLabel: model.LabelValue(server.PublicNet.IPv4.IP.String()),
	// 				},
	// 			},
	// 			Labels: model.LabelSet{
	// 				model.AddressLabel:                      model.LabelValue(server.PublicNet.IPv4.IP.String()),
	// 				model.LabelName(projectLabel):           model.LabelValue(project),
	// 				model.LabelName(nameLabel):              model.LabelValue(server.Name),
	// 				model.LabelName(statusLabel):            model.LabelValue(server.Status),
	// 				model.LabelName(publicIPv4Label):        model.LabelValue(server.PublicNet.IPv4.IP.String()),
	// 				model.LabelName(publicIPv6Label):        model.LabelValue(server.PublicNet.IPv6.IP.String()),
	// 				model.LabelName(serverTypeNameLabel):    model.LabelValue(server.ServerType.Name),
	// 				model.LabelName(serverTypeCoresLabel):   model.LabelValue(strconv.Itoa(int(server.ServerType.Cores))),
	// 				model.LabelName(serverTypeMemoryLabel):  model.LabelValue(strconv.Itoa(int(server.ServerType.Memory))),
	// 				model.LabelName(serverTypeDiskLabel):    model.LabelValue(strconv.Itoa(int(server.ServerType.Disk))),
	// 				model.LabelName(serverTypeStorageLabel): model.LabelValue(server.ServerType.StorageType),
	// 				model.LabelName(serverTypeCPULabel):     model.LabelValue(server.ServerType.CPUType),
	// 				model.LabelName(datacenterNameLabel):    model.LabelValue(server.Datacenter.Name),
	// 				model.LabelName(locationNameLabel):      model.LabelValue(server.Datacenter.Location.Name),
	// 				model.LabelName(locationCityLabel):      model.LabelValue(server.Datacenter.Location.City),
	// 				model.LabelName(locationCountryLabel):   model.LabelValue(server.Datacenter.Location.Country),
	// 				model.LabelName(imageTypeLabel):         model.LabelValue(imageType),
	// 				model.LabelName(imageNameLabel):         model.LabelValue(imageName),
	// 				model.LabelName(osFlavorLabel):          model.LabelValue(osFlavor),
	// 				model.LabelName(osVersionLabel):         model.LabelValue(osVersion),
	// 			},
	// 		}

	// 		for key, value := range server.Labels {
	// 			target.Labels[model.LabelName(labelPrefix+key)] = model.LabelValue(value)
	// 		}

	// 		level.Debug(d.logger).Log(
	// 			"msg", "Server added",
	// 			"project", project,
	// 			"source", target.Source,
	// 		)

	// 		current[target.Source] = struct{}{}
	// 		targets = append(targets, target)
	// 	}

	// }

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
