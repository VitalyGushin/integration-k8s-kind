// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration_k8s_kind_test

import (
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/edwarnicke/exechelper"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
)

type BasicTestsSuite struct {
	suite.Suite
	options []*exechelper.Option
}

func (s *BasicTestsSuite) TestDeployMemoryRegistry() {
	s.Require().NoError(exechelper.Run("kubectl apply -f ./deployments/memory-registry.yaml", s.options...))
	s.Require().NoError(exechelper.Run("kubectl delete -f ./deployments/memory-registry.yaml", s.options...))
}

func (s *BasicTestsSuite) SetupSuite() {
	writer := logrus.StandardLogger().Writer()
	s.options = []*exechelper.Option{
		exechelper.WithStderr(writer),
		exechelper.WithStdout(writer),
	}
	s.Require().NoError(exechelper.Run("kubectl apply -f deployments/spire/spire-namespace.yaml", s.options...))
	s.Require().NoError(exechelper.Run("kubectl apply -f deployments/spire", s.options...))
}

func (s *BasicTestsSuite) TestK8sClient() {
	path := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", path)
	s.NoError(err)
	_, err = kubernetes.NewForConfig(config)
	s.NoError(err)
}

func (s *BasicTestsSuite) TestDeployAlpine() {
	s.Require().NoError(exechelper.Run("kubectl delete -f ./deployments/alpine.yaml", s.options...))

	s.Require().NoError(exechelper.Run("kubectl apply -f ./deployments/alpine.yaml", s.options...))
	s.Require().NoError(exechelper.Run("kubectl wait --for=condition=ready pod -l app=alpine", s.options...))
}

func TestBasic(t *testing.T) {
	suite.Run(t, &BasicTestsSuite{})
}
