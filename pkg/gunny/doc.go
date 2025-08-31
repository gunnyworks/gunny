// Gunny facilitates rendering of static text content using templates and data
// from various sources. It attempts to follow the UNIX philosophy of being
// interoperable with other tooling as far as possible (e.g. in accepting data
// piped in via stdin).
//
// Potential use cases include static web site and configuration file
// generation at present.
//
// # Usage
//
// Usage of Gunny from code is demonstrated in the example below.
//
//	const rawTemplate = `Hello {{name}}!`
//
//	cliArgs := []string{"name=Michael"}
//
//	// Construct a rendering pipeline that uses an in-memory template. By
//	// default, if no output writer is specified, all output will go to
//	// os.Stdout.
//	pipeline, err := gunny.NewPipeline(
//		gunny.WithMustacheTemplateFromReader(strings.NewReader(rawTemplate)),
//		gunny.WithDataFromNameValuePairs(cliArgs),
//	)
//	// Handle error...
//
//	// Renders "Hello Michael!" to stdout.
//	err = pipeline.Render(context.Background())
//	// Handle error...
package gunny
