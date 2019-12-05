package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Galdoba/utils"
)

const (
	taskStatusDone  = "__done__"
	taskStatusError = "__error__"
)

/*
TZ:
1 прочитать задание
2 определить файлы которые надо ждать
3 найти нужные файлы (если их нет - ждать 2 часа с запросом каждую минуту)
4 проверить нужно ли резать видео


*/

var inputVideoDuration vidDuration

type vidDuration struct {
	frames       int
	timeCodePrem string
	timeCodeFF   string
}

func NewDuration(premiereCode string) *vidDuration {
	dur := vidDuration{}
	dur.timeCodePrem = premiereCode
	dur.timeCodeFF = premToFF(premiereCode)
	dur.frames = premToFrames(premiereCode)
	return &dur
}

type task struct {
	taskLineNum     int
	inputVideo      string
	intendentLenght string
	startTimeCode   string
	outputFileName  string
	taskStatus      string
}

func filesReady(t task) bool {
	fmt.Println("Uncut Input Video:", t.inputVideo)
	for !fileAvailable(t.inputVideo) {
		return false
		timeStamp()
		fmt.Println("WAIT filesready(t task)")
		time.Sleep(time.Second * 3)
	}
	//time.Sleep(time.Second * 30)
	fmt.Println("Input File Available")
	inputDuration := ffToPrem(videoDuration(t.inputVideo))
	fmt.Println("Input Video Duration:", inputDuration)
	fmt.Println("Intended Duration:", t.intendentLenght)
	if inputDuration == t.intendentLenght {
		renameFile(t.inputVideo, "cut_"+t.inputVideo)
	} else {
		cutVideo(t.startTimeCode, t.intendentLenght, t.inputVideo)
	}
	vid1Found := false
	aud1Found := false
	aud2Found := false
	audioFiles := predictAudioNames(t.outputFileName)
	for k := range audioFiles {
		fmt.Println("waiting for:", audioFiles[k])
	}
	needAudio := len(audioFiles)
	audio1 := ""
	audio2 := ""
	switch len(audioFiles) {
	case 1:
		audio1 = audioFiles[0]

	case 2:
		audio1 = audioFiles[0]
		audio2 = audioFiles[1]
	}
	//for (vid1Found && aud1Found && aud2Found) == false {

	timeStamp()
	vid1Found = fileAvailable("cut_" + t.inputVideo)
	fmt.Print("cut_" + t.inputVideo)
	if vid1Found {
		fmt.Println("...ok")
	} else {
		fmt.Println("...not ready")
	}
	fmt.Println("Need audio files:", needAudio)
	aud1Found = fileAvailable(audio1)
	if aud1Found {
		fmt.Println(audio1 + "...ok")
	}
	if needAudio == 2 {
		aud2Found = fileAvailable(audio2)
		if aud2Found {
			fmt.Println(audio2 + "...ok")
		}
	} else {
		aud2Found = true
		audio2 = "Not Using"
	}
	// if (vid1Found && aud1Found && aud2Found) == false {
	// 	time.Sleep(time.Second * 8)
	// }
	fmt.Println("")
	//}
	return vid1Found && aud1Found && aud2Found
}

func taskArgument(taskArgs []string, i int) string {
	if len(taskArgs) > i {
		return taskArgs[i]
	}
	return ""
}

func formTasks() []task {
	tasksLine := readTasks()
	var allTasks []task
	for i := range tasksLine {
		if tasksLine[i] == "" {
			continue
		}
		taskArgs := strings.Split(tasksLine[i], " ")
		currentTask := task{}
		currentTask.taskLineNum = i
		currentTask.inputVideo = taskArgument(taskArgs, 0)
		currentTask.intendentLenght = taskArgument(taskArgs, 1)
		currentTask.startTimeCode = taskArgument(taskArgs, 2)
		currentTask.outputFileName = taskArgument(taskArgs, 3)
		currentTask.taskStatus = taskArgument(taskArgs, 4)
		allTasks = append(allTasks, currentTask)
	}
	return allTasks
}

func main() {
	end := false
	for !end {
		allTasks := formTasks()
		if allDoneStatusOn(allTasks) {
			fmt.Println("All Tasks Done: nothing to do")
			break
		}
		for i := range allTasks {
			doTask(allTasks[i])
			fmt.Println("Tasks", i, "END")
		}
		if !allDoneStatusOn(allTasks) {
			fmt.Println("Wait because not all done")
			time.Sleep(time.Second * 3)
		}

	}
}

func allDoneStatusOn(allTasks []task) bool {
	end := true
	for i := range allTasks {
		if allTasks[i].taskStatus != taskStatusDone {
			return false
		}
	}
	return end
}

func doTask(t task) {

	if t.taskStatus == taskStatusError {
		fmt.Println("TASK STATUS:", t.taskStatus)
		return
	}
	if t.taskStatus == taskStatusDone {
		fmt.Println("TASK STATUS:", t.taskStatus)
		return
	}
	fmt.Println("Search Input Video:", t.inputVideo)
	if !isVideoAvailable(t) {
		fmt.Println("Input Video:", t.inputVideo, "___NOT FOUND")
		return
	}
	if !videoDurationConfirmed(t) {
		fmt.Println("Video duration:", t.inputVideo, "___NOT MATCH")
		return
	}
	audio1, audio2 := getAudioNames(t)
	if audio2 == "" {
		audio2 = audio1
	}
	if !allInputReady(t.inputVideoFile(), audio1, audio2) {
		fmt.Println("STATUS REPORT (files not ready)")
		return
	}

	//	panic("step 3")

	// for !fileAvailable(t.inputVideo) {
	// 	utils.EditLineInFile("FFTask.txt", t.taskLineNum, t.asString()+" "+"notFound")
	// 	return
	// 	timeStamp()
	// 	time.Sleep(time.Second * 3)
	// }
	//time.Sleep(time.Second * 30)
	t.changeTaskStatus(taskStatusDone)
	fmt.Println("valid task", t)
}

func isVideoAvailable(t task) bool {
	if !fileAvailable(t.inputVideo) {
		if !fileAvailable("cut_" + t.inputVideo) {
			return false
		}
		fmt.Println("Found cutted")
		//return false
	}
	return true
}

func (t *task) inputVideoFile() string {
	return "cut_" + t.inputVideo
}

func allInputReady(video, audio1, audio2 string) bool {
	if fileAvailable(video) && fileAvailable(audio1) && fileAvailable(audio2) {
		return true
	}
	return false
}

func getAudioNames(t task) (string, string) {
	audio1 := ""
	audio2 := ""
	audioFiles := predictAudioNames(t.outputFileName)
	for k := range audioFiles {
		fmt.Println("waiting for:", audioFiles[k])
	}
	//needAudio := len(audioFiles)
	switch len(audioFiles) {
	case 1:
		audio1 = audioFiles[0]

	case 2:
		audio1 = audioFiles[0]
		audio2 = audioFiles[1]
	}
	return audio1, audio2
}

func videoDurationConfirmed(t task) bool {
	inputDuration := ""
	inputDuration = ffToPrem(videoDuration(t.inputVideo))
	if fileAvailable("cut_" + t.inputVideo) {
		inputDuration = ffToPrem(videoDuration("cut_" + t.inputVideo))
	}
	if inputDuration != t.intendentLenght {
		fmt.Println("Cut&Create")
		cutVideo(t.startTimeCode, t.intendentLenght, t.inputVideo)
	} else {
		fmt.Println("rename")
		renameFile(t.inputVideo, "cut_"+t.inputVideo)
	}
	inputDuration = ffToPrem(videoDuration("cut_" + t.inputVideo))
	// if inputDuration != t.intendentLenght {
	// 	panic("UNKNOWN ERROR")
	// 	return false
	// }
	return true
}

func (t *task) asString() string {
	return t.inputVideo + " " + t.intendentLenght + " " + t.startTimeCode + " " + t.outputFileName
}

func (t *task) changeTaskStatus(newStatus string) {
	t.taskStatus = newStatus
	utils.EditLineInFile("FFTask.txt", t.taskLineNum, t.asString()+" "+newStatus)
}

func main0() {
	video1 := ""
	audio1 := ""
	audio2 := ""
	shortNme := ""
	needAudio := 0
	tasksLine := readTasks()
	for i := range tasksLine {
		fmt.Println(tasksLine[i])
	}

	allTasks := formTasks()
	fmt.Println("allTasks[i]")
	for i := range allTasks {
		fmt.Println(allTasks[i])
	}
	panic("Stop")
	for i := range tasksLine {
		if tasksLine[i] == "" {
			continue
		}
		fmt.Println("taskLine", i)
		///////////////////////////////////////////собираем информацию о задании
		taskArgs := strings.Split(tasksLine[i], " ")
		currentTask := task{}
		video1 = taskArgs[0]
		currentTask.inputVideo = taskArgs[0]
		shortNme = shortName(video1)
		fmt.Println(shortNme)
		//taskLen := taskArgs[1]
		currentTask.intendentLenght = taskArgs[1]
		//taskStart := taskArgs[2]
		currentTask.startTimeCode = taskArgs[2]
		intendedResult := taskArgs[3]
		currentTask.outputFileName = taskArgs[3]
		if len(taskArgs) > 4 {
			currentTask.taskStatus = taskArgs[4]
			if currentTask.taskStatus == taskStatusDone {
				fmt.Println("SKIP TASK")
				continue
			}
		}
		if !filesReady(currentTask) {
			fmt.Println("SKIP TASK files not ready")
			continue
		}
		audioFiles := predictAudioNames(currentTask.outputFileName)
		for k := range audioFiles {
			fmt.Println("waiting for:", audioFiles[k])
		}
		needAudio = len(audioFiles)
		switch len(audioFiles) {
		case 1:
			audio1 = audioFiles[0]

		case 2:
			audio1 = audioFiles[0]
			audio2 = audioFiles[1]
		}
		////////////////////////////////////////проверяеи и режем видео если надо (перенесено в проверку тасков)
		// fmt.Println("Uncut Input Video:", video1)
		// for !fileAvailable(video1) {
		// 	timeStamp()
		// 	time.Sleep(time.Second * 30)
		// }
		// //time.Sleep(time.Second * 30)
		// fmt.Println("Input File Available")
		// inputDuration := ffToPrem(videoDuration(video1))
		// fmt.Println("Input Video Duration:", inputDuration)
		// fmt.Println("Intended Duration:", taskLen)
		// if inputDuration == taskLen {
		// 	renameFile(video1, "cut_"+video1)
		// } else {
		// 	cutVideo(taskStart, taskLen, video1)
		// }
		//////////////////////////////////////собираем все нужные файлы
		// vid1Found := false
		// aud1Found := false
		// aud2Found := false
		// for (vid1Found && aud1Found && aud2Found) == false {

		// 	timeStamp()
		// 	vid1Found = fileAvailable("cut_" + video1)
		// 	fmt.Print("cut_" + video1)
		// 	if vid1Found {
		// 		fmt.Println("...ok")
		// 	} else {
		// 		fmt.Println("...not ready")
		// 	}
		// 	fmt.Println("Need audio files:", needAudio)
		// 	aud1Found = fileAvailable(audio1)
		// 	if aud1Found {
		// 		fmt.Println(audio1 + "...ok")
		// 	}
		// 	if needAudio == 2 {
		// 		aud2Found = fileAvailable(audio2)
		// 		if aud2Found {
		// 			fmt.Println(audio2 + "...ok")
		// 		}
		// 	} else {
		// 		aud2Found = true
		// 		audio2 = "Not Using"
		// 	}
		// 	if (vid1Found && aud1Found && aud2Found) == false {
		// 		time.Sleep(time.Second * 8)
		// 	}
		// 	fmt.Println("лишний цикл")
		// 	continue
		// }

		/////////////////////////////////////////Муксим

		/////////////////////////////////
		fmt.Println("Intended Result:", intendedResult)
		////////hd rus20
		if haveTag("cut_"+video1, "_hd_") && haveTag(audio1, "_rus20") && needAudio == 1 {
			runConsole("ffmpeg", "-i", "cut_"+video1, "-i", audio1, "-vcodec", "copy", "-acodec", "ac3", "-ab", "320k", "OUT_"+intendedResult+".mp4")

		}
		///////hd  rus51
		if haveTag("cut_"+video1, "_hd_") && haveTag(audio1, "_rus51") && needAudio == 1 {
			runConsole("ffmpeg", "-i", "cut_"+video1, "-i", audio1, "-vcodec", "copy", "-acodec", "ac3", "-ab", "640k", "OUT_"+intendedResult+".mp4")
			//ffmpeg -i "%f_fullname%" -i "%f_name%.aac" -vcodec copy -acodec ac3 -ab 640k "OUT_%f_name%_ar6%f_ext%"
		}
		///////sd  rus20
		if haveTag("cut_"+video1, "_sd_") && haveTag(audio1, "_rus20") && needAudio == 1 {
			runConsole("ffmpeg", "-i", "cut_"+video1, "-i", audio1, "-vcodec", "copy", "-acodec", "copy", "-ab", "320k", "OUT_"+intendedResult+".mp4")
			//ffmpeg -i "%f_fullname%" -i "%f_name%.aac" -vcodec copy -acodec ac3 -ab 640k "OUT_%f_name%_ar6%f_ext%"
		}
		////////sd rus20 Eng20
		if haveTag("cut_"+video1, "_sd_") && haveTag(audio1, "_rus20") && haveTag(audio2, "_eng20") && needAudio == 2 {
			runConsole("ffmpeg", "-i", audio1, "-acodec", "ac3", "-ab", "320k", shortNme+"_rus20.ac3")
			modAudio1 := shortNme + "_rus20.ac3"
			runConsole("ffmpeg", "-i", audio2, "-acodec", "ac3", "-ab", "320k", shortNme+"_eng20.ac3")
			modAudio2 := shortNme + "_eng20.ac3"
			//mkvmerge -o "OUT_%name%_ar2e2.mkv" -d 0 --language 0:rus --default-track 0:1 -A ="%vsrc%" -a 0 --language 0:rus --default-track 0:1 ="%name%_rus.%aext%" -a 0 --language 0:eng --default-track 0:0 ="%name%_eng.%aext%"
			runConsole("mkvmerge", "-o", "OUT_"+intendedResult+".mkv", "-d", "0", "--language", "0:rus", "--default-track", "0:1", "-A", "=cut_"+video1,
				"-a", "0", "--language", "0:rus", "--default-track", "0:1", "="+modAudio1,
				"-a", "0", "--language", "0:eng", "--default-track", "0:0", "="+modAudio2,
			)
		}
		////////hd rus20 Eng20
		if haveTag("cut_"+video1, "_hd_") && haveTag(audio1, "_rus20") && haveTag(audio2, "_eng20") && needAudio == 2 {
			fmt.Println(audio1)

			runConsole("ffmpeg", "-i", audio1, "-acodec", "ac3", "-ab", "320k", shortNme+"_rus20.ac3")
			modAudio1 := shortNme + "_rus20.ac3"
			fmt.Println(modAudio1)

			runConsole("ffmpeg", "-i", audio2, "-acodec", "ac3", "-ab", "320k", shortNme+"_eng20.ac3")
			modAudio2 := shortNme + "_eng20.ac3"
			//mkvmerge -o "OUT_%name%_ar2e2.mkv" -d 0 --language 0:rus --default-track 0:1 -A ="%vsrc%" -a 0 --language 0:rus --default-track 0:1 ="%name%_rus.%aext%" -a 0 --language 0:eng --default-track 0:0 ="%name%_eng.%aext%"
			runConsole("mkvmerge", "-o", "OUT_"+intendedResult+".mkv", "-d", "0", "--language", "0:rus", "--default-track", "0:1", "-A", "=cut_"+video1,
				"-a", "0", "--language", "0:rus", "--default-track", "0:1", "="+modAudio1,
				"-a", "0", "--language", "0:eng", "--default-track", "0:0", "="+modAudio2,
			)
		}
		////////hd rus20 Eng51
		if haveTag("cut_"+video1, "_hd_") && haveTag(audio1, "_rus20") && haveTag(audio2, "_eng51") && needAudio == 2 {
			runConsole("ffmpeg", "-i", audio1, "-acodec", "ac3", "-ab", "320k", shortNme+"_rus20.ac3")
			modAudio1 := shortNme + "_rus20.ac3"
			runConsole("ffmpeg", "-i", audio2, "-acodec", "ac3", "-ab", "640k", shortNme+"_eng51.ac3")
			modAudio2 := shortNme + "_eng51.ac3"
			//mkvmerge -o "OUT_%name%_ar2e2.mkv" -d 0 --language 0:rus --default-track 0:1 -A ="%vsrc%" -a 0 --language 0:rus --default-track 0:1 ="%name%_rus.%aext%" -a 0 --language 0:eng --default-track 0:0 ="%name%_eng.%aext%"
			runConsole("mkvmerge", "-o", "OUT_"+intendedResult+".mkv", "-d", "0", "--language", "0:rus", "--default-track", "0:1", "-A", "=cut_"+video1,
				"-a", "0", "--language", "0:rus", "--default-track", "0:1", "="+modAudio1,
				"-a", "0", "--language", "0:eng", "--default-track", "0:0", "="+modAudio2,
			)
		}
		////////hd rus51 Eng20
		if haveTag("cut_"+video1, "_hd_") && haveTag(audio1, "_rus51") && haveTag(audio2, "_eng20") && needAudio == 2 {
			runConsole("ffmpeg", "-i", audio1, "-acodec", "ac3", "-ab", "640k", shortNme+"_rus51.ac3")
			modAudio1 := shortNme + "_rus51.ac3"
			runConsole("ffmpeg", "-i", audio2, "-acodec", "ac3", "-ab", "320k", shortNme+"_eng20.ac3")
			modAudio2 := shortNme + "_eng20.ac3"
			//mkvmerge -o "OUT_%name%_ar2e2.mkv" -d 0 --language 0:rus --default-track 0:1 -A ="%vsrc%" -a 0 --language 0:rus --default-track 0:1 ="%name%_rus.%aext%" -a 0 --language 0:eng --default-track 0:0 ="%name%_eng.%aext%"
			runConsole("mkvmerge", "-o", "OUT_"+intendedResult+".mkv", "-d", "0", "--language", "0:rus", "--default-track", "0:1", "-A", "=cut_"+video1,
				"-a", "0", "--language", "0:rus", "--default-track", "0:1", "="+modAudio1,
				"-a", "0", "--language", "0:eng", "--default-track", "0:0", "="+modAudio2,
			)
		}
		////////hd rus51 Eng51
		if haveTag("cut_"+video1, "_hd_") && haveTag(audio1, "_rus51") && haveTag(audio2, "_eng51") && needAudio == 2 {
			runConsole("ffmpeg", "-i", audio1, "-acodec", "ac3", "-ab", "640k", shortNme+"_rus51.ac3")
			modAudio1 := shortNme + "_rus51.ac3"
			runConsole("ffmpeg", "-i", audio2, "-acodec", "ac3", "-ab", "640k", shortNme+"_eng51.ac3")
			modAudio2 := shortNme + "_eng51.ac3"
			//mkvmerge -o "OUT_%name%_ar2e2.mkv" -d 0 --language 0:rus --default-track 0:1 -A ="%vsrc%" -a 0 --language 0:rus --default-track 0:1 ="%name%_rus.%aext%" -a 0 --language 0:eng --default-track 0:0 ="%name%_eng.%aext%"
			runConsole("mkvmerge", "-o", "OUT_"+intendedResult+".mkv", "-d", "0", "--language", "0:rus", "--default-track", "0:1", "-A", "=cut_"+video1,
				"-a", "0", "--language", "0:rus", "--default-track", "0:1", "="+modAudio1,
				"-a", "0", "--language", "0:eng", "--default-track", "0:0", "="+modAudio2,
			)
		}

		fmt.Println("Task Complete")
		fmt.Println("")
		fmt.Println("TODO: CHANGE TASK LINE HERE")
		fmt.Println("taskLine", i)
		utils.EditLineInFile("FFTask.txt", i, tasksLine[i]+" "+taskStatusDone)
		fmt.Println("")
	}
	fmt.Println("Test for errors")
	waitEnter()
}

func catchTail(video1, audio1, audio2 string) string {
	tail := "ERROR"
	////////hd rus20
	if haveTag("cut_"+video1, "_hd_") && haveTag(audio1, "_rus20") && haveTag(audio2, "_rus20") {
		//runConsole("ffmpeg", "-i", "cut_"+video1, "-i", audio1, "-vcodec", "copy", "-acodec", "ac3", "-ab", "320k", "OUT_"+intendedResult+".mp4")
		tail = "hd_rus20_rus20"
	}
	return tail
}

func waitEnter() {
	fmt.Println("Press 'Enter'")
	input := bufio.NewScanner(os.Stdin) //Creating a Scanner that will read the input from the console
	for input.Scan() {
		if input.Text() == "" {
			break
		}
		fmt.Println(input.Text())
	}
}

func timeStamp() {
	t := time.Now()
	fmt.Println(t.Format("Time: 2006-01-02 15:04:05"))
}

func shortName(fullName string) string {
	validExtentions := []string{
		".mp4",
		".mov",
		".mpeg",
		".mpa",
		".ac3",
		".aac",
	}
	for i := range validExtentions {
		if strings.Contains(fullName, validExtentions[i]) {
			return strings.Split(fullName, validExtentions[i])[0]
		}
	}
	return ""
}

func haveTag(name, tag string) bool {
	have := strings.Contains(name, tag)
	return have
}

// func pingTask(task string) []string {

// 	fmt.Println("Current task:", task)
// 	taskArgs := strings.Split(task, " ")
// 	taskLen := taskArgs[1]
// 	taskStart := taskArgs[2]
// 	taskResult := taskArgs[3]
// 	taskInput := taskArgs[0]
// 	if taskLen == ffToPrem(videoDuration(taskInput)) {
// 		fmt.Println("No need to cut")
// 	} else {
// 		fmt.Println("Need to cut:", taskStart, "from start")
// 	}
// 	names := predictAudioNames(taskResult)
// 	return names
// }

func predictAudioNames(resultName string) []string {
	var expected []string
	//Pod_solncem_toskany_2003__hd_q0w1_ar2e2
	tagPool := strings.Split(resultName, "__")
	if len(tagPool) < 2 {
		fmt.Println(tagPool)
		return expected
	}
	allTags := tagPool[1]
	tags := strings.Split(allTags, "_")
	audioTag := tags[len(tags)-1]
	fmt.Println(audioTag)
	if strings.Contains(allTags, "sd_") {
		if strings.Contains(audioTag, "r2") {
			fmt.Println("case sd r2")
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_rus20.mpa"
			expected = append(expected, audName)
		}

		if strings.Contains(audioTag, "e2") {
			fmt.Println("case sd e2")
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_eng20.mpa"
			expected = append(expected, audName)
		}
	}

	if strings.Contains(allTags, "hd_") {
		if strings.Contains(audioTag, "r2") {
			fmt.Println("case hd r2")
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_rus20.aac"
			expected = append(expected, audName)
		}
		if strings.Contains(audioTag, "r6") {
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_rus51.aac"
			expected = append(expected, audName)
		}

		if strings.Contains(audioTag, "e2") {
			fmt.Println("case e2")
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_eng20.aac"
			expected = append(expected, audName)
		}

		if strings.Contains(audioTag, "e6") {
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_eng51.aac"
			expected = append(expected, audName)
		}
	}
	if strings.Contains(allTags, "3d_") {
		if strings.Contains(audioTag, "r2") {
			fmt.Println("case 3d r2")
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_rus20.aac"
			expected = append(expected, audName)
		}
		if strings.Contains(audioTag, "r6") {
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_rus51.aac"
			expected = append(expected, audName)
		}

		if strings.Contains(audioTag, "e2") {
			fmt.Println("case e2")
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_eng20.aac"
			expected = append(expected, audName)
		}

		if strings.Contains(audioTag, "e6") {
			audName := tagPool[0] + "_"
			for t := range tags {
				if t == len(tags)-1 {
					break
				}
				audName = audName + "_" + tags[t]
			}
			audName = audName + "_eng51.aac"
			expected = append(expected, audName)
		}
	}

	fmt.Println("done")

	return expected
}

func cutVideo(timeStart, timeLen, file string) string {
	outputFile := "cut_" + file
	if fileAvailable(file) {
		cmd := exec.Command("ffmpeg", "-i", file, "-map", "0:0", "-vcodec", "copy", "-an", "-t", premToFF(timeLen), "-ss", premToFF(timeStart), outputFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	} else {
		fmt.Println("Warning:", file, "is not available")
		fmt.Println("Waiting...")
		fmt.Println("")
	}
	return outputFile
}

func runConsole(program string, args ...string) {
	var line []string
	line = append(line, program)
	line = append(line, args...)
	fmt.Println("Run:", line)
	time.Sleep(time.Second)
	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// exists returns whether the given file or directory exists
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	fmt.Println(err)
	return true, err
}

func renameFile(old, new string) error {
	err := os.Rename(old, new)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func fileAvailable(file string) bool {
	err := os.Rename(file, "RENAMED_"+file)
	if err != nil {
		fmt.Println("file '" + file + "' is not available...")
		return false
	}
	os.Rename("RENAMED_"+file, file)
	return true
}

func readTasks() []string {
	var tasks []string
	file, err := os.Open("FFtask.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tasks = append(tasks, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tasks
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func videoDuration(file string) string {
	if fileAvailable(file) {
		cmd := exec.Command("ffmpeg", "-i", file)

		output, _ := cmd.CombinedOutput()
		stringOUT := string(output)
		str1 := strings.Split(stringOUT, "Duration: ")
		//time.Sleep(time.Second * 1)
		if len(str1) > 0 {
			durationSTR := strings.Split(str1[1], ", ")
			return durationSTR[0]
		}
	}
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Run()
	return "00:00:00.00"
}

func ffToPrem(duration string) string {
	premiereDur := ""
	timeArgs := strings.Split(duration, ":")
	hours, errHour := strconv.Atoi(timeArgs[0])
	if errHour != nil {
		panic(errHour)
	}
	minutes, errMin := strconv.Atoi(timeArgs[1])
	if errMin != nil {
		panic(errMin)
	}
	secs, err := strconv.ParseFloat(timeArgs[2], 64)
	if err != nil {
		panic(err)
	}
	secsInt := int(secs)
	strHour := strconv.Itoa(hours)
	if hours < 10 {
		strHour = "0" + strHour
	}
	strMin := strconv.Itoa(minutes)
	if minutes < 10 {
		strMin = "0" + strMin
	}
	strSec := strconv.Itoa(secsInt)
	if secsInt < 10 {
		strSec = "0" + strSec
	}
	part := secs
	for part > 1 {
		part = part - 1
	}
	part = toFixed(part, 3)
	part = part / 0.04
	frame := int(part)
	strFrame := strconv.Itoa(frame)
	if frame < 10 {
		strFrame = "0" + strFrame
	}
	if frame == 25 {
		strFrame = "00"
	}
	premiereDur = strHour + ":" + strMin + ":" + strSec + ":" + strFrame
	return premiereDur
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func premToFF(duration string) string {
	parts := strings.Split(duration, ":")
	secsInt, _ := strconv.Atoi(parts[2])
	partsInt, _ := strconv.Atoi(parts[3])
	partsFl := float64(partsInt)*40/1000 + float64(secsInt)
	sec := floatToString(partsFl)
	return parts[0] + ":" + parts[1] + ":" + sec
}

func premToFrames(duration string) int {
	parts := strings.Split(duration, ":")
	hour, _ := strconv.Atoi(parts[0])
	min, _ := strconv.Atoi(parts[1])
	sec, _ := strconv.Atoi(parts[2])
	frm, _ := strconv.Atoi(parts[3])
	return frm + 25*sec + 1500*min + 90000*hour
}

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done() // Need to signal to waitgroup that this goroutine is done
}

/*
mux_hd_ruseng_names
set CURRENTDIR=%CD%
set cwd=%~dp0
for /f "delims="  %%i IN ( 'type names.txt' ) DO (
call :process "%%i"
)

pause
exit /b 0

:process
set f_fullname=%~1
set f_path=%~p1
set f_name=%~n1
set f_ext=%~x1

rem ffmpeg -i %f_name%_rus.mpa -map 0:0 -acodec ac3 -ab 320k %f_name%_rus.ac3
rem ffmpeg -i %f_name%_eng.mpa -map 0:0 -acodec ac3 -ab 320k %f_name%_eng.ac3

mkvmerge -o "%f_name%_ar2e2.mkv" -d 0 --language 0:rus --default-track 0:1 -A ="d:\Work\petr_proj\__done\OUT\%f_name%.mp4" -a 0 --language 0:rus --default-track 0:1 ="%f_name%_rus.mpa" -a 0 --language 0:eng --default-track 0:0 ="%f_name%_eng.mpa"

rem -s 0 --language 0:rus --default-track 0:0 ="%f_name%.srt"

exit /b 0




mux_mutevideo_1audio_names_51
set CURRENTDIR=%CD%
set cwd=%~dp0

for /f "delims="  %%i IN ( 'type names.txt' ) DO (
call :process "%%i"
)

pause
exit /b 0

:process
set f_fullname=%~1
set f_path=%~p1
set f_name=%~n1
set f_ext=%~x1
ffmpeg -i "%f_fullname%" -i "%f_name%.aac" -vcodec copy -acodec ac3 -ab 640k "OUT_%f_name%_ar6%f_ext%"
exit /b 0



mux_mutevideo_1audio_names_20
set CURRENTDIR=%CD%
set cwd=%~dp0

for /f "delims="  %%i IN ( 'type names.txt' ) DO (
call :process "%%i"
)

pause
exit /b 0

:process
set f_fullname=%~1
set f_path=%~p1
set f_name=%~n1
set f_ext=%~x1
ffmpeg -i "%f_fullname%" -i "%f_name%.aac" -vcodec copy -acodec ac3 -ab 320k "OUT_%f_name%_ar2%f_ext%"
exit /b 0




mux_mutevideo_1audio_names_mpa
set CURRENTDIR=%CD%
set cwd=%~dp0

for /f "delims="  %%i IN ( 'type names.txt' ) DO (
call :process "%%i"
)

pause
exit /b 0

:process
set f_fullname=%~1
set f_path=%~p1
set f_name=%~n1
set f_ext=%~x1
ffmpeg -i "%f_fullname%" -i "%f_name%.mpa" -vcodec copy -acodec copy -ab 320k "OUT_%f_name%_ar2%f_ext%"
exit /b 0


mutevideo_audio_engrus5ch
set aext=ac3
set name=shest_dney_sem_nochey_1998__hd_q0w0
set vsrc=f:\Work\petr_proj\__done\OUT\shest_dney_sem_nochey_1998__hd_q0w0.mp4

ffmpeg -i "%name%_rus.aac" -acodec ac3 -ab 640k %name%_rus.%aext%
ffmpeg -i "%name%_eng.aac" -acodec ac3 -ab 640k %name%_eng.%aext%

mkvmerge -o "%name%_ar6e6.mkv" -d 0 --language 0:rus --default-track 0:1 -A ="%vsrc%" -a 0 --language 0:rus --default-track 0:1 ="%name%_rus.%aext%" -a 0 --language 0:eng --default-track 0:0 ="%name%_eng.%aext%"

rem  -s 0 --language 0:rus --default-track 0:0 ="%1.srt"

pause



*/
