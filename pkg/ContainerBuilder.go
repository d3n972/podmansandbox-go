package pkg

import (
	"crypto/md5"
	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/martinlindhe/base36"
	"github.com/opencontainers/runtime-spec/specs-go"
	"time"
)

type MountType string
type MountOpt uint16

const (
	MountOptRO MountOpt = iota + 4
)
const (
	MountTypeBind   MountType = "bind"
	MountTypeVolume MountType = "volume"
	MountTypeTmpFS  MountType = "tmpfs"
	MountTypeImage  MountType = "image"
	MountTypeDevPts MountType = "devpts"
)

type ContainerBuilder struct {
	gen *specgen.SpecGenerator
}

func NewContainerBuilder(image string) *ContainerBuilder {
	spec := specgen.NewSpecGenerator(image, false)
	spec.Labels = map[string]string{}
	return &ContainerBuilder{
		gen: spec,
	}
}
func (c *ContainerBuilder) GenName() *ContainerBuilder {
	c.gen.Name = base36.EncodeBytes(
		md5.New().
			Sum(
				[]byte(
					time.Now().Format(time.RFC3339),
				),
			),
	)
	return c
}
func (c *ContainerBuilder) SetName(name string) *ContainerBuilder {
	c.gen.Name = name
	return c
}
func (c *ContainerBuilder) MapPort(outer uint16, inner uint16) *ContainerBuilder {
	c.gen.PortMappings = append(c.gen.PortMappings, types.PortMapping{
		ContainerPort: inner,
		HostPort:      outer,
	})
	return c
}
func (c *ContainerBuilder) Mount(t MountType, src string, dst string) *ContainerBuilder {
	c.gen.Mounts = append(c.gen.Mounts, specs.Mount{
		Destination: dst,
		Type:        string(t),
		Source:      src,
	})
	return c
}
func (c *ContainerBuilder) Label(key string, value string) *ContainerBuilder {
	c.gen.Labels[key] = value
	return c
}

func (c *ContainerBuilder) BuildSpecification() *specgen.SpecGenerator {
	return c.gen
}
