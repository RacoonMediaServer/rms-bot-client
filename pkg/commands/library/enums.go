package library

import (
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
)

type listName string
type movieType string

const (
	listNameFavourites listName = "Избранное"
	listNameWatchList  listName = "Онлайн-просмотр"
	listNameArchive    listName = "Архив"
)

var listNameToMovieType = map[listName]rms_library.List{
	listNameFavourites: rms_library.List_Favourites,
	listNameWatchList:  rms_library.List_WatchList,
	listNameArchive:    rms_library.List_Archive,
}

func getListName(movieType rms_library.List) listName {
	for name, t := range listNameToMovieType {
		if t == movieType {
			return name
		}
	}
	return ""
}

func getListType(name listName) (rms_library.List, bool) {
	t, ok := listNameToMovieType[name]
	return t, ok
}
