// Package watch implements file-system watching for the vaultpull config file.
//
// When the config changes on disk the registered Handler is invoked after a
// short debounce delay so that rapid successive saves are collapsed into a
// single sync operation.
//
// Usage:
//
//	w := watch.New(".vaultpull.yaml", 200*time.Millisecond, syncFunc)
//	if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package watch
