package taskrun

import (
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"sort"
)

func SortByNamespace(trs []v1alpha1.TaskRun) {
	sort.Sort(byNamespace(trs))
}

type byNamespace []v1alpha1.TaskRun

func (prs byNamespace) compareNamespace(ins, jns string) (lt, eq bool) {
	lt, eq = ins < jns, ins == jns
	return lt, eq
}

func (trs byNamespace) Len() int      { return len(trs) }
func (trs byNamespace) Swap(i, j int) { trs[i], trs[j] = trs[j], trs[i] }
func (trs byNamespace) Less(i, j int) bool {
	var lt, eq bool
	if lt, eq = trs.compareNamespace(trs[i].Namespace, trs[j].Namespace); eq {
		if trs[j].Status.StartTime == nil {
			return false
		}
		if trs[i].Status.StartTime == nil {
			return true
		}
		return trs[j].Status.StartTime.Before(trs[i].Status.StartTime)
	}
	return lt
}
