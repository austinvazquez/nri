/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package generate_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	rspec "github.com/opencontainers/runtime-spec/specs-go"
	rgen "github.com/opencontainers/runtime-tools/generate"

	"github.com/containerd/nri/pkg/api"
	xgen "github.com/containerd/nri/pkg/runtime-tools/generate"
)

func TestGenerate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generate Suite")
}

var _ = Describe("Adjustment", func() {
	When("nil", func() {
		It("does not modify the Spec", func() {
			var (
				spec   = makeSpec()
				adjust *api.ContainerAdjustment
			)

			rg := &rgen.Generator{Config: spec}
			xg := xgen.SpecGenerator(rg)

			Expect(xg).ToNot(BeNil())
			Expect(xg.Adjust(adjust)).To(Succeed())
			Expect(spec).To(Equal(makeSpec()))
		})
	})

	When("empty", func() {
		It("does not modify the Spec", func() {
			var (
				spec   = makeSpec()
				adjust = &api.ContainerAdjustment{
					Linux: &api.LinuxContainerAdjustment{},
				}
			)

			rg := &rgen.Generator{Config: spec}
			xg := xgen.SpecGenerator(rg)

			Expect(xg).ToNot(BeNil())
			Expect(xg.Adjust(adjust)).To(Succeed())
			Expect(spec).To(Equal(makeSpec()))
		})
	})

	When("has memory limit", func() {
		It("adjusts Spec correctly", func() {
			var (
				spec   = makeSpec()
				adjust = &api.ContainerAdjustment{
					Linux: &api.LinuxContainerAdjustment{
						Resources: &api.LinuxResources{
							Memory: &api.LinuxMemory{
								Limit: api.Int64(11111),
							},
						},
					},
				}
			)

			rg := &rgen.Generator{Config: spec}
			xg := xgen.SpecGenerator(rg)

			Expect(xg).ToNot(BeNil())
			Expect(xg.Adjust(adjust)).To(Succeed())
			Expect(spec).To(Equal(makeSpec(withMemoryLimit(11111), withMemorySwap(11111))))
		})
	})

	When("has CPU shares", func() {
		It("adjusts Spec correctly", func() {
			var (
				spec   = makeSpec()
				adjust = &api.ContainerAdjustment{
					Linux: &api.LinuxContainerAdjustment{
						Resources: &api.LinuxResources{
							Cpu: &api.LinuxCPU{
								Shares: api.UInt64(11111),
							},
						},
					},
				}
			)

			rg := &rgen.Generator{Config: spec}
			xg := xgen.SpecGenerator(rg)

			Expect(xg).ToNot(BeNil())
			Expect(xg.Adjust(adjust)).To(Succeed())
			Expect(spec).To(Equal(makeSpec(withCPUShares(11111))))
		})
	})

	When("has CPU quota", func() {
		It("adjusts Spec correctly", func() {
			var (
				spec   = makeSpec()
				adjust = &api.ContainerAdjustment{
					Linux: &api.LinuxContainerAdjustment{
						Resources: &api.LinuxResources{
							Cpu: &api.LinuxCPU{
								Quota: api.Int64(11111),
							},
						},
					},
				}
			)

			rg := &rgen.Generator{Config: spec}
			xg := xgen.SpecGenerator(rg)

			Expect(xg).ToNot(BeNil())
			Expect(xg.Adjust(adjust)).To(Succeed())
			Expect(spec).To(Equal(makeSpec(withCPUQuota(11111))))
		})
	})

	When("has CPU period", func() {
		It("adjusts Spec correctly", func() {
			var (
				spec   = makeSpec()
				adjust = &api.ContainerAdjustment{
					Linux: &api.LinuxContainerAdjustment{
						Resources: &api.LinuxResources{
							Cpu: &api.LinuxCPU{
								Period: api.UInt64(11111),
							},
						},
					},
				}
			)

			rg := &rgen.Generator{Config: spec}
			xg := xgen.SpecGenerator(rg)

			Expect(xg).ToNot(BeNil())
			Expect(xg.Adjust(adjust)).To(Succeed())
			Expect(spec).To(Equal(makeSpec(withCPUPeriod(11111))))
		})
	})

	When("has cpuset CPUs", func() {
		It("adjusts Spec correctly", func() {
			var (
				spec   = makeSpec()
				adjust = &api.ContainerAdjustment{
					Linux: &api.LinuxContainerAdjustment{
						Resources: &api.LinuxResources{
							Cpu: &api.LinuxCPU{
								Cpus: "5,6",
							},
						},
					},
				}
			)

			rg := &rgen.Generator{Config: spec}
			xg := xgen.SpecGenerator(rg)

			Expect(xg).ToNot(BeNil())
			Expect(xg.Adjust(adjust)).To(Succeed())
			Expect(spec).To(Equal(makeSpec(withCPUSetCPUs("5,6"))))
		})
	})

	When("has cpuset mems", func() {
		It("adjusts Spec correctly", func() {
			var (
				spec   = makeSpec()
				adjust = &api.ContainerAdjustment{
					Linux: &api.LinuxContainerAdjustment{
						Resources: &api.LinuxResources{
							Cpu: &api.LinuxCPU{
								Mems: "5,6",
							},
						},
					},
				}
			)

			rg := &rgen.Generator{Config: spec}
			xg := xgen.SpecGenerator(rg)

			Expect(xg).ToNot(BeNil())
			Expect(xg.Adjust(adjust)).To(Succeed())
			Expect(spec).To(Equal(makeSpec(withCPUSetMems("5,6"))))
		})
	})
})

type specOption func(*rspec.Spec)

func withMemoryLimit(v int64) specOption {
	return func(spec *rspec.Spec) {
		if spec.Linux == nil {
			spec.Linux = &rspec.Linux{}
		}
		if spec.Linux.Resources == nil {
			spec.Linux.Resources = &rspec.LinuxResources{}
		}
		if spec.Linux.Resources.Memory == nil {
			spec.Linux.Resources.Memory = &rspec.LinuxMemory{}
		}
		spec.Linux.Resources.Memory.Limit = &v
	}
}

func withMemorySwap(v int64) specOption {
	return func(spec *rspec.Spec) {
		if spec.Linux == nil {
			spec.Linux = &rspec.Linux{}
		}
		if spec.Linux.Resources == nil {
			spec.Linux.Resources = &rspec.LinuxResources{}
		}
		if spec.Linux.Resources.Memory == nil {
			spec.Linux.Resources.Memory = &rspec.LinuxMemory{}
		}
		spec.Linux.Resources.Memory.Swap = &v
	}
}

func withCPUShares(v uint64) specOption {
	return func(spec *rspec.Spec) {
		if spec.Linux == nil {
			spec.Linux = &rspec.Linux{}
		}
		if spec.Linux.Resources == nil {
			spec.Linux.Resources = &rspec.LinuxResources{}
		}
		if spec.Linux.Resources.CPU == nil {
			spec.Linux.Resources.CPU = &rspec.LinuxCPU{}
		}
		spec.Linux.Resources.CPU.Shares = &v
	}
}

func withCPUQuota(v int64) specOption {
	return func(spec *rspec.Spec) {
		if spec.Linux == nil {
			spec.Linux = &rspec.Linux{}
		}
		if spec.Linux.Resources == nil {
			spec.Linux.Resources = &rspec.LinuxResources{}
		}
		if spec.Linux.Resources.CPU == nil {
			spec.Linux.Resources.CPU = &rspec.LinuxCPU{}
		}
		spec.Linux.Resources.CPU.Quota = &v
	}
}

func withCPUPeriod(v uint64) specOption {
	return func(spec *rspec.Spec) {
		if spec.Linux == nil {
			spec.Linux = &rspec.Linux{}
		}
		if spec.Linux.Resources == nil {
			spec.Linux.Resources = &rspec.LinuxResources{}
		}
		if spec.Linux.Resources.CPU == nil {
			spec.Linux.Resources.CPU = &rspec.LinuxCPU{}
		}
		spec.Linux.Resources.CPU.Period = &v
	}
}

func withCPUSetCPUs(v string) specOption {
	return func(spec *rspec.Spec) {
		if spec.Linux == nil {
			spec.Linux = &rspec.Linux{}
		}
		if spec.Linux.Resources == nil {
			spec.Linux.Resources = &rspec.LinuxResources{}
		}
		if spec.Linux.Resources.CPU == nil {
			spec.Linux.Resources.CPU = &rspec.LinuxCPU{}
		}
		spec.Linux.Resources.CPU.Cpus = v
	}
}

func withCPUSetMems(v string) specOption {
	return func(spec *rspec.Spec) {
		if spec.Linux == nil {
			spec.Linux = &rspec.Linux{}
		}
		if spec.Linux.Resources == nil {
			spec.Linux.Resources = &rspec.LinuxResources{}
		}
		if spec.Linux.Resources.CPU == nil {
			spec.Linux.Resources.CPU = &rspec.LinuxCPU{}
		}
		spec.Linux.Resources.CPU.Mems = v
	}
}

func makeSpec(options ...specOption) *rspec.Spec {
	spec := &rspec.Spec{
		Process: &rspec.Process{},
		Linux: &rspec.Linux{
			Resources: &rspec.LinuxResources{
				Memory: &rspec.LinuxMemory{
					Limit: Int64(12345),
				},
				CPU: &rspec.LinuxCPU{
					Shares: Uint64(45678),
					Quota:  Int64(87654),
					Period: Uint64(54321),
					Cpus:   "0-111",
					Mems:   "0-4",
				},
			},
		},
	}
	for _, o := range options {
		o(spec)
	}
	return spec
}

func Int64(v int64) *int64 {
	return &v
}

func Uint64(v uint64) *uint64 {
	return &v
}
