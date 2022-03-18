package reconciler

import (
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/sirupsen/logrus"
)

var logger = common.Logger().WithFields(logrus.Fields{
	"pkg": "controlplane/reconciler",
})
