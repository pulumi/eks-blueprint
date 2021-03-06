// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const providerName = "ssp"
const version = "0.0.1"

type module struct {
	version semver.Version
}

func (m *module) Version() semver.Version {
	return m.version
}

func (m *module) Construct(ctx *pulumi.Context, name, typ, urn string) (r pulumi.Resource, err error) {
	switch typ {
	case "ssp:index:SharedServicesPlatform":
		r = &Ssp{}
	default:
		return nil, fmt.Errorf("unknown resource type: %s", typ)
	}

	err = ctx.RegisterResource(typ, name, nil, r, pulumi.URN_(urn))
	return
}

// Serve launches the gRPC server for the resource provider.
func Serve(providerName, version string, schema []byte) {
	// Start gRPC service.
	pulumi.RegisterResourceModule("ssp", "index", &module{version: semver.MustParse(version)})

	if err := provider.MainWithOptions(provider.Options{
		Name:      providerName,
		Version:   version,
		Schema:    schema,
		Construct: construct,
		Call:      call,
	}); err != nil {
		cmdutil.ExitError(err.Error())
	}
}
