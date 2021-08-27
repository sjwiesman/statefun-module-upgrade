package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "statefun-module-upgrade",
		Short: "Convert Apache Flink Stateful Function module.yaml to >= 3.1 format",
		Run: func(cmd *cobra.Command, args []string) {
			convert()
		},
	}

	input string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&input, "input", "", "input file (default stdin)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error: %s", err)
	}
}

func convert() {
	b, err := readModule()
	if err != nil {
		log.Panicf("failed to read input: %s", err)
	}

	var module Module
	err = yaml.Unmarshal(b, &module)

	d, err := module.MarshalText()
	fmt.Printf(string(d))
}

func readModule() ([]byte, error) {
	var err error
	r := os.Stdin
	if len(input) != 0 {
		r, err = os.Open(input)
		defer func(r *os.File) {
			_ = r.Close()
		}(r)
		if err != nil {
			panic(err)
		}
	}
	return ioutil.ReadAll(r)
}

type Endpoint struct {
	Kind string `yaml:"kind"`
	Spec struct {
		Functions       string                 `yaml:"functions"`
		UrlPathTemplate string                 `yaml:"urlPathTemplate"`
		Timeouts        map[string]interface{} `yaml:"timeouts"`
	} `yaml:"spec"`
}

func (e *Endpoint) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var legacy struct {
		Endpoint struct {
			Meta struct {
				Kind string `yaml:"kind"`
			} `yaml:"meta"`
			Spec struct {
				Functions       string                 `yaml:"functions"`
				UrlPathTemplate string                 `yaml:"urlPathTemplate"`
				Timeouts        map[string]interface{} `yaml:"timeouts"`
			} `yaml:"spec"`
		} `yaml:"endpoint"`
	}

	if err := unmarshal(&legacy); err != nil {
		return fmt.Errorf("failed to unmarshal legacy endpoint: %w", err)
	}

	(*e).Kind = "io.statefun.endpoints.v2/http"
	(*e).Spec.Functions = legacy.Endpoint.Spec.Functions
	(*e).Spec.UrlPathTemplate = legacy.Endpoint.Spec.UrlPathTemplate
	(*e).Spec.Timeouts = legacy.Endpoint.Spec.Timeouts

	return nil
}

func (e *Endpoint) MarshalYAML() (interface{}, error) {
	return struct {
		Kind string `yaml:"kind"`
		Spec struct {
			Functions       string                 `yaml:"functions"`
			UrlPathTemplate string                 `yaml:"urlPathTemplate"`
			Timeouts        map[string]interface{} `yaml:"timeouts"`
		} `yaml:"spec"`
	}{
		Kind: e.Kind,
		Spec: e.Spec,
	}, nil
}

type Ingress struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"spec"`
}

type KafkaIngressSpec struct {
	Id              string      `yaml:"id"`
	Address         string      `yaml:"address"`
	ConsumerGroupId string      `yaml:"consumerGroupId,omitempty"`
	StartupPosition interface{} `yaml:"startupPosition,omitempty"`
	Properties      interface{} `yaml:"properties,omitempty"`
	Topics          interface{} `yaml:"topics"`
}

type KinesisIngressSpec struct {
	Id                     string      `yaml:"id"`
	AwsRegion              interface{} `yaml:"awsRegion,omitempty"`
	AwsCredentials         interface{} `yaml:"awsCredentials,omitempty"`
	StartupPosition        interface{} `yaml:"startupPosition,omitempty"`
	Streams                interface{} `yaml:"streams,omitempty"`
	ClientConfigProperties interface{} `yaml:"clientConfigProperties,omitempty"`
}

func (i *Ingress) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var t struct {
		Ingress struct {
			Meta struct {
				Type string `yaml:"type"`
			}
		} `yaml:"ingress"`
	}

	if err := unmarshal(&t); err != nil {
		return fmt.Errorf("failed to unmarshal legacy ingress: %w", err)
	}

	switch t.Ingress.Meta.Type {
	case "io.statefun.kafka/ingress":
		var legacy struct {
			Ingress struct {
				Meta struct {
					Type string `yaml:"type"`
					Id   string `yaml:"id"`
				}
				Spec KafkaIngressSpec `yaml:"spec"`
			} `yaml:"ingress"`
		}

		if err := unmarshal(&legacy); err != nil {
			return fmt.Errorf("failed to unmarshal legacy kafka ingress: %w", err)
		}
		(*i).Kind = "io.statefun.kafka.v1/ingress"
		legacy.Ingress.Spec.Id = legacy.Ingress.Meta.Id
		(*i).Spec = legacy.Ingress.Spec
	case "io.statefun.kinesis/ingress":
		var legacy struct {
			Ingress struct {
				Meta struct {
					Type string `yaml:"type"`
					Id   string `yaml:"id"`
				}
				Spec KinesisIngressSpec `yaml:"spec"`
			} `yaml:"ingress"`
		}

		if err := unmarshal(&legacy); err != nil {
			return fmt.Errorf("failed to unmarshal legacy kinesis ingress: %w", err)
		}
		(*i).Kind = "io.statefun.kinesis.v1/ingress"
		legacy.Ingress.Spec.Id = legacy.Ingress.Meta.Id
		(*i).Spec = legacy.Ingress.Spec
	default:
		return fmt.Errorf("unknown ingress type %s", t.Ingress.Meta.Type)
	}

	return nil
}

func (i *Ingress) MarshalYAML() (interface{}, error) {
	return struct {
		Kind string      `yaml:"kind"`
		Spec interface{} `yaml:"spec"`
	}{
		Kind: i.Kind,
		Spec: i.Spec,
	}, nil
}

type Egress struct {
	Kind string      `yaml:"kind"`
	Spec interface{} `yaml:"spec"`
}

type KafkaEgressSpec struct {
	Id               string      `yaml:"id"`
	Address          string      `yaml:"address"`
	DeliverySemantic interface{} `yaml:"deliverySemantic,omitempty"`
	Properties       interface{} `yaml:"properties,omitempty"`
}

type KinesisEgressSpec struct {
	Id                     string      `yaml:"id"`
	AwsRegion              interface{} `yaml:"awsRegion,omitempty"`
	AwsCredentials         interface{} `yaml:"awsCredentials,omitempty"`
	MaxOutstandingRecords  interface{} `yaml:"maxOutstandingRecords,omitempty"`
	ClientConfigProperties interface{} `yaml:"clientConfigProperties,omitempty"`
}

func (e *Egress) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var t struct {
		Ingress struct {
			Meta struct {
				Type string `yaml:"type"`
			}
		} `yaml:"egress"`
	}

	if err := unmarshal(&t); err != nil {
		return fmt.Errorf("failed to unmarshal legacy egress: %w", err)
	}

	switch t.Ingress.Meta.Type {
	case "io.statefun.kafka/egress":
		var legacy struct {
			Egress struct {
				Meta struct {
					Type string `yaml:"type"`
					Id   string `yaml:"id"`
				}
				Spec KafkaEgressSpec `yaml:"spec"`
			} `yaml:"egress"`
		}

		if err := unmarshal(&legacy); err != nil {
			return fmt.Errorf("failed to unmarshal legacy kafka egress: %w", err)
		}
		(*e).Kind = "io.statefun.kafka.v1/egress"
		legacy.Egress.Spec.Id = legacy.Egress.Meta.Id
		(*e).Spec = legacy.Egress.Spec
	case "io.statefun.kinesis/ingress":
		var legacy struct {
			Egress struct {
				Meta struct {
					Type string `yaml:"type"`
					Id   string `yaml:"id"`
				}
				Spec KinesisEgressSpec `yaml:"spec"`
			} `yaml:"egress"`
		}

		if err := unmarshal(&legacy); err != nil {
			return fmt.Errorf("failed to unmarshal legacy kinesis egress: %w", err)
		}
		(*e).Kind = "io.statefun.kinesis.v1/egress"
		legacy.Egress.Spec.Id = legacy.Egress.Meta.Id
		(*e).Spec = legacy.Egress.Spec
	default:
		return fmt.Errorf("unknown ingress type %s", t.Ingress.Meta.Type)
	}

	return nil
}

func (e *Egress) MarshalYAML() (interface{}, error) {
	return struct {
		Kind string      `yaml:"kind"`
		Spec interface{} `yaml:"spec"`
	}{
		Kind: e.Kind,
		Spec: e.Spec,
	}, nil
}

type Module []interface{}

func (m *Module) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var legacy struct {
		Module struct {
			Spec struct {
				Endpoints []struct {
					Endpoint `yaml:"endpoint"`
				} `yaml:"endpoints"`
				Ingresses []struct {
					Ingress `yaml:"ingress"`
				} `yaml:"ingresses"`
				Egresses []struct {
					Egress `yaml:"egress"`
				} `yaml:"egresses"`
			} `yaml:"spec"`
		} `yaml:"module"`
	}

	if err := unmarshal(&legacy); err != nil {
		return fmt.Errorf("failed to unmarshal module.yaml: %w", err)
	}

	components := make([]interface{}, 0)

	for _, component := range legacy.Module.Spec.Endpoints {
		components = append(components, &component)
	}
	for _, component := range legacy.Module.Spec.Ingresses {
		components = append(components, &component)
	}
	for _, component := range legacy.Module.Spec.Egresses {
		components = append(components, &component)
	}

	*m = components

	return nil
}

func (m *Module) MarshalText() (text []byte, err error) {
	builder := strings.Builder{}
	for _, component := range *m {
		d, err := yaml.Marshal(component)
		if err != nil {
			return nil, err
		}

		builder.Write(d)
		builder.WriteString("---\n")
	}

	text = []byte(builder.String())
	return
}
