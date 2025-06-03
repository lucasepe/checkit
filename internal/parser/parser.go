package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Parse(src io.Reader) (CheckList, error) {
	scanner := bufio.NewScanner(src)

	var (
		list         CheckList
		currentGroup *Group
		currentItem  *Item
		titleSet     bool
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // skip empty lines
		}

		switch {
		case strings.HasPrefix(line, "# "):
			if titleSet {
				// skip multiple titles
				continue
			}

			list.Title = strings.TrimSpace(line[2:])
			titleSet = true

		case strings.HasPrefix(line, "## "):
			groupTitle := strings.TrimSpace(line[3:])
			group := Group{Title: groupTitle}
			list.Groups = append(list.Groups, group)
			currentGroup = &list.Groups[len(list.Groups)-1]
			currentItem = nil

		case strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* "):
			if currentGroup == nil {
				return list, fmt.Errorf("item does not belong to any group: %q", line)
			}
			itemTitle := strings.TrimSpace(line[2:])
			item := Item{Title: itemTitle}
			currentGroup.Items = append(currentGroup.Items, item)
			currentItem = &currentGroup.Items[len(currentGroup.Items)-1]

		case strings.HasPrefix(line, ">"):
			if currentItem == nil {
				return list, fmt.Errorf("note does not belong to any item: %q", line)
			}
			note := strings.TrimSpace(line[1:])
			currentItem.Notes = append(currentItem.Notes, note)

		default:
			// skip all the rest
			continue
		}
	}

	err := scanner.Err()

	return list, err
}
