package privacy

import (
	"strings"

	containerName "github.com/google/go-containerregistry/pkg/name"
)

var allowedBuildpackSources = []struct {
	registry, repositoryPrefix string
}{
	// Paketo
	{
		registry:         "gcr.io",
		repositoryPrefix: "paketo-buildpacks/",
	}, {
		registry:         "index.docker.io",
		repositoryPrefix: "paketobuildpacks/",
	},
	// Google Buildpacks
	{
		registry:         "gcr.io",
		repositoryPrefix: "buildpacks/",
	},
	// Heroku
	{
		registry:         "public.ecr.aws",
		repositoryPrefix: "heroku-buildpacks/",
	},
}

func FilterBuildpacks(buildpacks []string) []string {
	result := make([]string, 0, len(buildpacks))
	for _, buildpack := range buildpacks {
		ref, err := containerName.ParseReference(strings.ToLower(buildpack))
		if err != nil {
			result = append(result, "<error>")
			continue
		}

		registry := ref.Context().Registry.Name()
		repository := ref.Context().RepositoryStr()

		allowed := false
		for _, allowedBuildpackSource := range allowedBuildpackSources {
			if registry == allowedBuildpackSource.registry && strings.HasPrefix(repository, allowedBuildpackSource.repositoryPrefix) {
				allowed = true
				break
			}
		}

		if allowed {
			result = append(result, buildpack)
		} else {
			result = append(result, "<redacted>")
		}
	}
	return result
}

var allowedEnvKeys = map[string]interface{}{"BP_JVM_VERSION": nil, "BP_NODE_VERSION": nil}

func FilterEnv(in map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for key, value := range in {
		_, allowed := allowedEnvKeys[key]
		if allowed {
			out[key] = value
		}
	}
	return out
}
