package installer_utils

import "fmt"

func getMessage(cmd CommandType, state string) (msg string) {
	switch cmd.Type {
	case COMMAND_INSTALL_TYPE:
		switch state {
		case DOING_STATE:
			msg = fmt.Sprintf("%s installing", cmd.Name)
			return
		case FAILED_STATE:
			msg = fmt.Sprintf("%s install failed", cmd.Name)
			return
		case DONE_STATE:
			msg = fmt.Sprintf("%s installed", cmd.Name)
			return
		}
	case COMMAND_IMPORT_KEY_TYPE:
		switch state {
		case DOING_STATE:
			msg = fmt.Sprintf("%s importing", cmd.Name)
			return
		case FAILED_STATE:
			msg = fmt.Sprintf("%s import failed", cmd.Name)
			return
		case DONE_STATE:
			msg = fmt.Sprintf("%s imported", cmd.Name)
			return
		}
	case COMMAND_UPGRADE_TYPE:
		switch state {
		case DOING_STATE:
			msg = fmt.Sprintf("%s upgrading", cmd.Name)
			return
		case FAILED_STATE:
			msg = fmt.Sprintf("%s upgrade failed", cmd.Name)
			return
		case DONE_STATE:
			msg = fmt.Sprintf("%s upgraded", cmd.Name)
			return
		}
	case COMMAND_PURGE_TYPE:
		switch state {
		case DOING_STATE:
			msg = fmt.Sprintf("%s purging", cmd.Name)
			return
		case FAILED_STATE:
			msg = fmt.Sprintf("%s purge failed", cmd.Name)
			return
		case DONE_STATE:
			msg = fmt.Sprintf("%s purged", cmd.Name)
			return
		}

	case COMMAND_UPDATE_TYPE:
		switch state {
		case DOING_STATE:
			msg = fmt.Sprintf("%s updating", cmd.Name)
			return
		case FAILED_STATE:
			msg = fmt.Sprintf("%s update failed", cmd.Name)
			return
		case DONE_STATE:
			msg = fmt.Sprintf("%s updated", cmd.Name)
			return
		}

	case COMMAND_EXEC_TYPE:
		switch state {
		case DOING_STATE:
			msg = fmt.Sprintf("%s executing", cmd.Name)
			return
		case FAILED_STATE:
			msg = fmt.Sprintf("%s execute failed", cmd.Name)
			return
		case DONE_STATE:
			msg = fmt.Sprintf("%s executed", cmd.Name)
			return
		}
	}
	return
}

func getDownloadMessage(dl DownloadType, state string) (msg string) {
	switch state {
	case DOING_STATE:
		msg = fmt.Sprintf("%s downloaing", dl.Name)
		return
	case FAILED_STATE:
		msg = fmt.Sprintf("%s download failed", dl.Name)
		return
	case DONE_STATE:
		msg = fmt.Sprintf("%s downloaded", dl.Name)
		return
	}
	return
}
