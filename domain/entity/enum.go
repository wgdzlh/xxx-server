package entity

const (
	RES_TYPE_IMG = "raster"
	RES_TYPE_VEC = "vector"

	SUB_TYPE_RET = "RESULT"
	SUB_TYPE_DOM = "DOM"

	STORE_TYPE_FLIGHT = "FLIGHT"

	STAT_OP_NUM          = 1
	STAT_OP_AREA         = 2
	STAT_OP_NUM_AND_AREA = 3

	STAT_TYPE_UP  = "update"
	STAT_TYPE_DEL = "delete"
	STAT_TYPE_ADD = "add"

	TASK_INIT    = "NotStarted"
	TASK_SKIPPED = "Skipped"
	TASK_IN_PROC = "InProc"
	TASK_STAGED  = "Staged"
	TASK_ST_PROC = "StageProc"
	TASK_DONE    = "Done"
	TASK_FAILED  = "Failed"

	TASK_TYPE_SHORT  = "short"
	TASK_TYPE_MID    = "mid"
	TASK_TYPE_FACTOR = "factor"

	TASK_GENRE_CRON = "cron"
	TASK_GENRE_DOWN = "download"
)
