package command

import (
	"context"
	"errors"

	"github.com/promhippie/prometheus-vcd-sd/pkg/action"
	"github.com/promhippie/prometheus-vcd-sd/pkg/config"
	"github.com/urfave/cli/v3"
)

// Server provides the sub-command to start the server.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: ServerFlags(cfg),
		Action: func(_ context.Context, cmd *cli.Command) error {
			logger := setupLogger(cfg)

			if cmd.IsSet("vcd.config") {
				if err := readConfig(cmd.String("vcd.config"), cfg); err != nil {
					logger.Error("Failed to read config",
						"err", err,
					)

					return err
				}
			}

			if cfg.Target.File == "" {
				logger.Error("Missing path for output.file")
				return errors.New("missing path for output.file")
			}

			if cmd.IsSet("vcd.url") && cmd.IsSet("vcd.username") && cmd.IsSet("vcd.password") && cmd.IsSet("vcd.org") && cmd.IsSet("vcd.vdc") {
				credentials := config.Credential{
					Project:  "default",
					URL:      cmd.String("vcd.url"),
					Insecure: cmd.Bool("vcd.insecure"),
					Username: cmd.String("vcd.username"),
					Password: cmd.String("vcd.password"),
					Org:      cmd.String("vcd.org"),
					Vdc:      cmd.String("vcd.vdc"),
				}

				cfg.Target.Credentials = append(
					cfg.Target.Credentials,
					credentials,
				)

				if credentials.URL == "" {
					logger.Error("Missing required vcd.url")
					return errors.New("missing required vcd.url")
				}

				if credentials.Username == "" {
					logger.Error("Missing required vcd.username")
					return errors.New("missing required vcd.username")
				}

				if credentials.Password == "" {
					logger.Error("Missing required vcd.password")
					return errors.New("missing required vcd.password")
				}

				if credentials.Org == "" {
					logger.Error("Missing required vcd.org")
					return errors.New("missing required vcd.org")
				}

				if credentials.Vdc == "" {
					logger.Error("Missing required vcd.vdc")
					return errors.New("missing required vcd.vdc")
				}
			}

			if len(cfg.Target.Credentials) == 0 {
				logger.Error("Missing any credentials")
				return errors.New("missing any credentials")
			}

			return action.Server(cfg, logger)
		},
	}
}

// ServerFlags defines the available server flags.
func ServerFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "web.address",
			Value:       "0.0.0.0:9000",
			Usage:       "Address to bind the metrics server",
			Sources:     cli.EnvVars("PROMETHEUS_VCD_WEB_ADDRESS"),
			Destination: &cfg.Server.Addr,
		},
		&cli.StringFlag{
			Name:        "web.path",
			Value:       "/metrics",
			Usage:       "Path to bind the metrics server",
			Sources:     cli.EnvVars("PROMETHEUS_VCD_WEB_PATH"),
			Destination: &cfg.Server.Path,
		},
		&cli.StringFlag{
			Name:        "web.config",
			Value:       "",
			Usage:       "Path to web-config file",
			Sources:     cli.EnvVars("PROMETHEUS_VCD_WEB_CONFIG"),
			Destination: &cfg.Server.Web,
		},
		&cli.StringFlag{
			Name:        "output.engine",
			Value:       "file",
			Usage:       "Enabled engine like file or http",
			Sources:     cli.EnvVars("PROMETHEUS_VCD_OUTPUT_ENGINE"),
			Destination: &cfg.Target.Engine,
		},
		&cli.StringFlag{
			Name:        "output.file",
			Value:       "/etc/prometheus/vcd.json",
			Usage:       "Path to write the file_sd config",
			Sources:     cli.EnvVars("PROMETHEUS_VCD_OUTPUT_FILE"),
			Destination: &cfg.Target.File,
		},
		&cli.IntFlag{
			Name:        "output.refresh",
			Value:       30,
			Usage:       "Discovery refresh interval in seconds",
			Sources:     cli.EnvVars("PROMETHEUS_VCD_OUTPUT_REFRESH"),
			Destination: &cfg.Target.Refresh,
		},
		&cli.StringFlag{
			Name:    "vcd.url",
			Value:   "",
			Usage:   "URL for the vCloud Director API",
			Sources: cli.EnvVars("PROMETHEUS_VCD_URL"),
		},
		&cli.BoolFlag{
			Name:    "vcd.insecure",
			Value:   false,
			Usage:   "Accept self-signed certs for the vCloud Director API",
			Sources: cli.EnvVars("PROMETHEUS_VCD_INSECURE"),
		},
		&cli.StringFlag{
			Name:    "vcd.username",
			Value:   "",
			Usage:   "Username for the vCloud Director API",
			Sources: cli.EnvVars("PROMETHEUS_VCD_USERNAME"),
		},
		&cli.StringFlag{
			Name:    "vcd.password",
			Value:   "",
			Usage:   "Password for the vCloud Director API",
			Sources: cli.EnvVars("PROMETHEUS_VCD_PASSWORD"),
		},
		&cli.StringFlag{
			Name:    "vcd.org",
			Value:   "",
			Usage:   "Organization for the vCloud Director API",
			Sources: cli.EnvVars("PROMETHEUS_VCD_ORG"),
		},
		&cli.StringFlag{
			Name:    "vcd.vdc",
			Value:   "",
			Usage:   "vDatacenter for the vCloud Director API",
			Sources: cli.EnvVars("PROMETHEUS_VCD_VDC"),
		},
		&cli.StringFlag{
			Name:    "vcd.config",
			Value:   "",
			Usage:   "Path to vCloud Director configuration file",
			Sources: cli.EnvVars("PROMETHEUS_VCD_CONFIG"),
		},
	}
}
