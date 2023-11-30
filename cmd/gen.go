package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/utils/fmtutil"
	"github.com/starudream/go-lib/tablew/v2"

	"github.com/starudream/aliyunxiao-generator/api"
	"github.com/starudream/aliyunxiao-generator/utils"
)

var genCmd = cobra.NewCommand(func(c *cobra.Command) {
	c.Use = "gen"

	var (
		now              = utils.NewTime()
		startWeekDateStr = now.StartOfWeek().Layout(time.DateOnly)
		endWeekDateStr   = now.EndOfWeek().Layout(time.DateOnly)
		startDateStr     = c.PersistentFlags().String("start-date", startWeekDateStr, "start date")
		endDataStr       = c.PersistentFlags().String("end-date", endWeekDateStr, "end date")

		weekly     = c.PersistentFlags().BoolP("weekly", "w", false, "weekly report")
		lastWeekly = c.PersistentFlags().BoolP("last-weekly", "W", false, "last weekly report")

		startDate, endDate utils.Carbon
	)

	c.RunE = func(cmd *cobra.Command, args []string) (err error) {
		switch {
		case *weekly:
			startDate, endDate = now.StartOfWeek(), now.EndOfWeek()
		case *lastWeekly:
			startDate, endDate = now.StartOfWeek().SubWeek(), now.EndOfWeek().SubWeek()
		default:
			startDate, err = utils.ParseTime(*startDateStr, time.DateOnly)
			if err != nil {
				return fmt.Errorf("start date is invalid: %w", err)
			}
			endDate, err = utils.ParseTime(*endDataStr, time.DateOnly)
			if err != nil {
				return fmt.Errorf("end date is invalid: %w", err)
			}
		}

		confirm := fmtutil.Scan(fmt.Sprintf("generate report from %q to %q, press `y` to continue: ", startDate.Layout(time.DateOnly), endDate.Layout(time.DateOnly)))
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			return nil
		}

		var (
			wim = map[string]*api.WorkItem{} // work item map (identifier -> work item)
			wtm = map[string]int{}           // work item time map (identifier -> hours)
			spm = map[string]string{}        // space map (identifier -> name)
			stm = map[string]int{}           // space time map (name -> hours)
			ths = 0                          // total hours
		)

		for t := startDate; t.Lte(endDate); t = t.AddDay() {
			items, err := api.ListWorkItem(t.ToDateString())
			if err != nil {
				return err
			}

			identifiers := make([]string, 0)
			for _, t := range items {
				identifiers = append(identifiers, t.Identifier)
			}
			identifiers = lo.Uniq(identifiers)

			hours, err := api.ListWorkItemTime(t.ToDateString(), identifiers)
			if err != nil {
				return err
			}

			for i := 0; i < len(items); i++ {
				t := items[i]

				if _, ok := wim[t.Identifier]; !ok {
					wim[t.Identifier] = items[i]
				}

				if spn, ok := spm[t.SpaceIdentifier]; ok {
					items[i].SpaceName = spn
				} else {
					sp, err := api.GetSpace(t.SpaceIdentifier)
					if err != nil {
						return err
					}
					items[i].SpaceName, spm[t.SpaceIdentifier] = sp.Name, sp.Name
				}

				hour := hours[t.Identifier]
				wtm[t.Identifier] += hour
				stm[t.SpaceName] += hour
				ths += hour
			}
		}

		its := make([][]string, 0) // table items

		for _, t := range wim {
			its = append(its, []string{t.SpaceName, t.WorkitemType.DisplayName, t.Subject, strconv.Itoa(wtm[t.Identifier])})
		}

		sort.Slice(its, func(i, j int) bool {
			if its[i][0] == its[j][0] {
				if its[i][1] == its[j][1] {
					return its[i][2] < its[j][2]
				}
				return its[i][1] < its[j][1]
			}
			return its[i][0] < its[j][0]
		})

		{

			table := tablew.Render(func(tw *tablew.Table) {
				tw.SetRowLine(true)
				tw.SetAutoWrapText(false)
				tw.SetAutoMergeCells(true)
				tw.SetColMinWidth(2, 100)
				tw.SetColumnAlignment([]int{tablew.ALIGN_CENTER, tablew.ALIGN_CENTER, tablew.ALIGN_CENTER, tablew.ALIGN_CENTER})
				tw.SetHeader([]string{"space", "category", "subject", "hours"})
				for i := 0; i < len(its); i++ {
					tw.Append(its[i])
				}
				tw.SetFooter([]string{"", "", "total", strconv.Itoa(ths)})
			})
			fmt.Printf("table:\n\n%s\n", table)
		}

		{
			buf := &bytes.Buffer{}
			space := ""
			for i := 0; i < len(its); i++ {
				if its[i][0] != space {
					space = its[i][0]
					buf.WriteString("\n")
					buf.WriteString("# ")
					buf.WriteString(space)
					buf.WriteString(" - ")
					buf.WriteString(strconv.Itoa(stm[space]))
					buf.WriteString("h\n")
				}
				buf.WriteString(" - ")
				buf.WriteString(its[i][2])
				buf.WriteString(" - ")
				buf.WriteString(its[i][3])
				buf.WriteString("h\n")
			}
			fmt.Printf("text:\n%s\n", buf.String())
		}

		return nil
	}
})

func init() {
	rootCmd.AddCommand(genCmd)
}
