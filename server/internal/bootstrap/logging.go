/**
 * dox
 * Copyright (C) 2026  OpenDox
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * @File    : logging.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-05-31
 * @Modified: 2026-05-31
 */

package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sharedlogging "github.com/opendox/dox/packages/shared/logging"
	"github.com/opendox/dox/packages/shared/version"
)

func loggingResourceFromSnapshot(snapshot *SettingSnapshot) sharedlogging.Resource {
	identity := snapshot.Setting.Identity
	serviceVersion := version.GetInfo().Version
	if override := strings.TrimSpace(snapshot.Setting.Logging.Resource.ServiceVersion); override != "" {
		serviceVersion = override
	}

	return sharedlogging.Resource{
		ServiceNamespace:      identity.Service.Namespace,
		ServiceName:           identity.Service.Name,
		ServiceInstanceID:     identity.Service.InstanceID,
		ServiceVersion:        serviceVersion,
		DeploymentEnvironment: string(identity.Deployment.Env),
		CloudRegion:           identity.Deployment.Region,
		CloudAvailabilityZone: identity.Deployment.Zone,
		K8sClusterName:        identity.Deployment.Cluster,
		K8sNamespaceName:      identity.Deployment.K8sNamespace,
		DoxOrganization:       identity.Organization.Name,
		DoxApplication:        identity.Application.Name,
		DoxRuntime:            string(identity.System.Runtime),
	}
}

func renderLoggingConfig(config sharedlogging.Config, resource sharedlogging.Resource) sharedlogging.Config {
	render := loggingPathRenderer(resource)

	config.Zap.OutputPaths = renderPaths(config.Zap.OutputPaths, render)
	config.Zap.ErrorOutputPaths = renderPaths(config.Zap.ErrorOutputPaths, render)
	if len(config.Cores) > 0 {
		cores := make([]sharedlogging.CoreConfig, len(config.Cores))
		copy(cores, config.Cores)
		for index := range cores {
			cores[index].OutputPaths = renderPaths(cores[index].OutputPaths, render)
		}
		config.Cores = cores
	}
	return config
}

func loggingPathRenderer(resource sharedlogging.Resource) *strings.Replacer {
	return strings.NewReplacer(
		"${service.namespace}", resource.ServiceNamespace,
		"${service.name}", resource.ServiceName,
		"${service.instance.id}", resource.ServiceInstanceID,
		"${service.version}", resource.ServiceVersion,
		"${deployment.environment.name}", resource.DeploymentEnvironment,
		"${cloud.region}", resource.CloudRegion,
		"${cloud.availability_zone}", resource.CloudAvailabilityZone,
		"${k8s.cluster.name}", resource.K8sClusterName,
		"${k8s.namespace.name}", resource.K8sNamespaceName,
		"${dox.organization}", resource.DoxOrganization,
		"${dox.application}", resource.DoxApplication,
		"${dox.runtime}", resource.DoxRuntime,
	)
}

func renderPaths(paths []string, render *strings.Replacer) []string {
	if paths == nil {
		return nil
	}
	rendered := make([]string, len(paths))
	for index, path := range paths {
		rendered[index] = render.Replace(path)
	}
	return rendered
}

func prepareLoggingPaths(config sharedlogging.Config) error {
	for _, path := range config.Zap.OutputPaths {
		if err := prepareLogFilePath(path); err != nil {
			return fmt.Errorf("prepare zap output path %q: %w", path, err)
		}
	}
	for _, path := range config.Zap.ErrorOutputPaths {
		if err := prepareLogFilePath(path); err != nil {
			return fmt.Errorf("prepare zap error output path %q: %w", path, err)
		}
	}
	for _, core := range config.Cores {
		if core.Type != sharedlogging.CoreTypeFile {
			continue
		}
		for _, path := range core.OutputPaths {
			if err := prepareLogFilePath(path); err != nil {
				return fmt.Errorf("prepare logging core %q output path %q: %w", core.Name, path, err)
			}
		}
	}
	return nil
}

func prepareLogFilePath(path string) error {
	path = strings.TrimSpace(path)
	if !isLocalLogFilePath(path) {
		return nil
	}

	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func isLocalLogFilePath(path string) bool {
	if path == "" {
		return false
	}
	switch path {
	case "stdout", "stderr":
		return false
	}
	return !strings.Contains(path, "://")
}
