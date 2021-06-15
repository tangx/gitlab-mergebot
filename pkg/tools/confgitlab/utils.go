package confgitlab

func IsMerged(mr *MergeRequest) bool {
	return mr.State == "merged"
}

func IsVaildMRCandidate(mr *MergeRequest) bool {
	if IsMerged(mr) {
		return false
	}

	if mr.MergeStatus != "can_be_merged" {
		return false
	}

	return true
}
