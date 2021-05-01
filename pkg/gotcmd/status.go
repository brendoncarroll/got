package gotcmd

import (
	"fmt"
	"io"

	"github.com/brendoncarroll/got/pkg/gotfs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(catCmd)
	rootCmd.AddCommand(checkCmd)
}

var statusCmd = &cobra.Command{
	Use:     "status",
	PreRunE: loadRepo,
	RunE: func(cmd *cobra.Command, args []string) error {
		w := cmd.OutOrStdout()
		name, _, err := repo.GetActiveVolume(ctx)
		if err != nil {
			return err
		}
		delta, err := repo.StagingDiff(ctx)
		if err != nil {
			return err
		}
		additions, err := delta.ListAdditionPaths(ctx, repo.StagingStore())
		if err != nil {
			return err
		}
		deletions, err := delta.ListDeletionPaths(ctx, repo.StagingStore())
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "ACTIVE: %s\n", name)

		printToBeCommitted(w, additions, deletions)

		fmt.Fprintf(w, "Changes not staged for commit:\n")
		fmt.Fprintf(w, "  (use \"got add <file>...\" to update what will be commited)\n")
		fmt.Fprintf(w, "  (use \"got clobber <file>...\" to discard changes in working directory)\n")
		// TODO: list paths with staged versions
		fmt.Fprintln(w, "    < TODO >")

		fmt.Fprintf(w, "Untracked files:\n")
		fmt.Fprintf(w, "  (use \"got add <file>...\" to include what will be commited)\n")
		// TODO: list paths not in staging
		fmt.Fprintln(w, "    < TODO >")

		return nil
	},
}

func printToBeCommitted(w io.Writer, additions, deletions []string) {
	if len(additions) == 0 && len(deletions) == 0 {
		return
	}
	fmt.Fprintf(w, "Changes to be committed:\n")
	for _, p := range additions {
		fmt.Fprintf(w, "    modified: %s\n", p)
	}
	for _, p := range deletions {
		fmt.Fprintf(w, "    deleted: %s\n", p)
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "lists the children of path in the current volume",
	PreRunE: loadRepo,
	RunE: func(cmd *cobra.Command, args []string) error {
		var p string
		if len(args) > 0 {
			p = args[0]
		}
		w := cmd.OutOrStdout()
		return repo.Ls(ctx, p, func(ent gotfs.DirEnt) error {
			_, err := fmt.Fprintf(w, "%v %s\n", ent.Mode, ent.Name)
			return err
		})
	},
}

var catCmd = &cobra.Command{
	Use:     "cat",
	Short:   "writes the contents of path in the current volume to stdout",
	PreRunE: loadRepo,
	RunE: func(cmd *cobra.Command, args []string) error {
		var p string
		if len(args) > 0 {
			p = args[0]
		}
		w := cmd.OutOrStdout()
		return repo.Cat(ctx, p, w)
	},
}

var checkCmd = &cobra.Command{
	Use:     "check",
	Short:   "checks contents of the current volume",
	PreRunE: loadRepo,
	RunE: func(cmd *cobra.Command, args []string) error {
		return repo.Check(ctx)
	},
}
