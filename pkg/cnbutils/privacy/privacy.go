package privacy

import (
	"strings"

	containerName "github.com/google/go-containerregistry/pkg/name"
)

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

		if registry == "gcr.io" && strings.HasPrefix(repository, "paketo-buildpacks/") {
			result = append(result, buildpack)
		} else if registry == "index.docker.io" && strings.HasPrefix(repository, "paketobuildpacks/") {
			result = append(result, buildpack)
		} else if registry == "gcr.io" && strings.HasPrefix(repository, "buildpacks/") {
			result = append(result, buildpack)
		} else if registry == "public.ecr.aws" && strings.HasPrefix(repository, "heroku-buildpacks/") {
			result = append(result, buildpack)
		} else {
			result = append(result, "<retracted>")
		}
	}
	return result
}
