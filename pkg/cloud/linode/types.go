package linode

import (
	"github.com/linode/linodego"
)

const logFeature = "cloudprovider/linode"

const (
	// linodeType4GB defines a Linode with 4GB ram, 80GB storage
	linodeType4GB = "g6-standard-2" // GET https://api.linode.com/v4/linode/types
	// regionFrankfurt defines a region in central Europe
	regionFrankfurt = "eu-central" // GET https://api.linode.com/v4/regions
)

type provider struct {
	client linodego.Client
}

type pollTestFn func() (ready bool, err error)
