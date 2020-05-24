package wrap

// Option is an option for modifying the actions of a Wrapper.
type Option func(Wrapper)

// AddBuildTag adds build tags to the generated files.
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
