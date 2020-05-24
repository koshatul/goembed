package wrap

type Option func(Wrapper)

func AddBuildTag(buildTags []string) Option {
	return func(w Wrapper) {
		switch v := w.(type) {
		case *AferoWrapper:
			v.buildTags = buildTags
		case *NoDepWrapper:
			v.buildTags = buildTags
		}
	}
}
