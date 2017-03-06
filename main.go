package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	git2go "gopkg.in/libgit2/git2go.v25"
)

func main() {
	inputPath := "."
	if len(os.Args) > 1 && os.Args[1] != "" {
		inputPath = os.Args[1]
	}
	outputPath := "output.csv"
	if len(os.Args) > 2 && os.Args[2] != "" {
		outputPath = os.Args[2]
	}

	out, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	parseGitHistory(inputPath, out)
}

func parseGitHistory(inputPath string, out io.Writer) {
	log.Println("Opening repo")
	repo, err := git2go.OpenRepository(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	rv, err := repo.Walk()
	rv.Sorting(git2go.SortTime)
	rv.PushHead()

	log.Println("Parsing each commit in history")
	head, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(out)
	if err := w.Write(csvHeader); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	oid := head.Target()
	for {
		commit, err := gitCommitFromOid(repo, oid)
		if err != nil {
			log.Fatal(err)
		}

		row := commitToCSVRow(commit)
		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}

		// Go to next commit
		if err := rv.Next(oid); err != nil {
			break
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	log.Println("Done.")
}

func gitCommitFromOid(repo *git2go.Repository, oid *git2go.Oid) (*git2go.Commit, error) {
	obj, err := repo.Lookup(oid)
	if err != nil {
		return nil, err
	}
	return obj.AsCommit()
}

var csvHeader = []string{
	"oid",
	"message",
	"time",
	"committer_name",
	"committer_email",
	"committer_time",
	"author_name",
	"author_email",
	"author_time",
	"parent_count",
}

func commitToCSVRow(c *git2go.Commit) []string {
	return []string{
		c.Object.Id().String(),      // oid
		c.Message(),                 // message
		c.Committer().When.String(), // time
		c.Committer().Name,          // committer_name
		c.Committer().Email,         // committer_email
		c.Committer().When.String(), // committer_time
		c.Author().Name,             // author_name
		c.Author().Email,            // author_email
		c.Author().When.String(),    // author_time

		strconv.FormatUint(uint64(c.ParentCount()), 10), // parent_count
	}
}
