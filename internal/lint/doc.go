// Package lint validates vaultpull configuration files and reports
// misconfigurations as structured issues with severity levels.
//
// Usage:
//
//	issues := lint.Run(cfg)
//	if lint.HasErrors(issues) {
//		for _, i := range issues {
//			fmt.Fprintln(os.Stderr, i)
//		}
//		os.Exit(1)
//	}
package lint
