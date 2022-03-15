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
	"time"

	"github.com/pkg/errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
)

type Abc struct {
	Name string `pulumi:"name"`
}

func construct(ctx *pulumi.Context, typ, name string, inputs provider.ConstructInputs,
	options pulumi.ResourceOption) (*provider.ConstructResult, error) {
	switch typ {
	case "ssp:index:SharedServicesPlatform":
		return constructSsp(ctx, name, inputs, options)
	default:
		return nil, errors.Errorf("unknown resource type %s", typ)
	}
}

func call(ctx *pulumi.Context, call string, inputs provider.CallArgs) (*provider.CallResult, error) {
	switch call {
	case "ssp:index:SharedServicesPlatform/addManagedNodeGroup":
		args := &AddManagedNodeGroupArgs{}

		_, err := inputs.CopyTo(args)
		if err != nil {
			return nil, err
		}

		memoCache := getMemoCache()
		resource, found := memoCache.Storage.Get("ssp")
		if !found {
			return nil, errors.New("Couldn't get SSP from cache")
		}

		ssp, ok := resource.(*Ssp)
		if !ok {
			return nil, errors.New("invalid resource type")
		}

		managedNodeGroup, err := ssp.addManagedNodeGroup(ctx, args)
		if err != nil {
			return nil, err
		}

		return pulumiprovider.NewCallResult(managedNodeGroup)

	case "ssp:index:SharedServicesPlatform/clusterAddon":
		args := &AddClusterAddOnArgs{}

		_, err := inputs.CopyTo(args)
		if err != nil {
			return nil, err
		}

		memoCache := getMemoCache()
		resource, found := memoCache.Storage.Get("ssp")
		if !found {
			return nil, errors.New("Couldn't get SSP from cache")
		}

		ssp, ok := resource.(*Ssp)
		if !ok {
			return nil, errors.New("invalid resource type")
		}

		h, err := ssp.clusterAddon(ctx, args)
		if err != nil {
			return nil, err
		}

		return pulumiprovider.NewCallResult(h)

	case "ssp:index:SharedServicesPlatform/addFargateProfile":
		args := &AddFargateProfileArgs{}

		_, err := inputs.CopyTo(args)
		if err != nil {
			return nil, err
		}

		memoCache := getMemoCache()
		resource, found := memoCache.Storage.Get("ssp")
		if !found {
			return nil, errors.New("Couldn't get SSP from cache")
		}

		ssp, ok := resource.(*Ssp)
		if !ok {
			return nil, errors.New("invalid resource type")
		}

		h, err := ssp.addFargateProfile(ctx, args)
		if err != nil {
			return nil, err
		}

		return pulumiprovider.NewCallResult(h)

	case "ssp:index:SharedServicesPlatform/onboardTeam":
		args := &OnboardTeamArgs{}

		_, err := inputs.CopyTo(args)
		if err != nil {
			return nil, err
		}

		memoCache := getMemoCache()
		resource, found := memoCache.Storage.Get("ssp")
		if !found {
			return nil, errors.New("Couldn't get SSP from cache")
		}

		ssp, ok := resource.(*Ssp)
		if !ok {
			return nil, errors.New("invalid resource type")
		}

		h, err := ssp.onboardTeam(ctx, args)
		if err != nil {
			return nil, err
		}

		return pulumiprovider.NewCallResult(h)

	case "ssp:index:SharedServicesPlatform/addNodeGroup":
	default:
		return nil, errors.Errorf("unknown call type %s", call)
	}

	return nil, errors.New("Shouldn't get here")
}

// constructSsp
func constructSsp(ctx *pulumi.Context, name string, inputs pulumiprovider.ConstructInputs,
	options pulumi.ResourceOption) (*pulumiprovider.ConstructResult, error) {

	args := &ClusterArgs{}
	if err := inputs.CopyTo(args); err != nil {
		return nil, errors.Wrap(err, "setting args")
	}

	memoCache = getMemoCache()

	component, err := NewSharedServicesPlatform(ctx, name, args, options)
	if err != nil {
		return nil, errors.Wrap(err, "creating component")
	}

	// should be keyed with the URN or name, but for now we'll just use ssp as I don't
	// know how to get the URN./name out in the call ^^
	err = memoCache.Storage.Add("ssp", component, time.Hour*1)
	if err != nil {
		return nil, errors.Wrap(err, "caching ssp")
	}

	return pulumiprovider.NewConstructResult(component)
}
