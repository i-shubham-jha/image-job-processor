package job

import (
	"fmt"
	"math/rand/v2"
	"retail_pulse/internal/files"
	"retail_pulse/internal/model"
	"retail_pulse/internal/service"
	"retail_pulse/internal/store"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// assumes that storesVisit has been validated by the caller
// and this id is marked as ongoing in db
func ProcessJob(id primitive.ObjectID, sv model.StoresVisit) {
	sm, err := store.NewStoreManager()

	if err != nil {
		fmt.Println(err)
		return
	}

	svs := service.NewStoresVisitService()

	for visitIndex, store := range sv.Visits {

		if !sm.StoreIDExists(store.StoreID) {
			svs.UpdateStoresVisitStatus(id, "failed", "store ID does not exist", store.StoreID)
			return
		}

		new_image_uuids := make([]string, len(store.ImageURLs))
		new_image_perims := make([]int64, len(store.ImageURLs))

		copy(new_image_uuids, store.ImageUUIDs)
		copy(new_image_perims, store.Perimeters)

		for i, img_url := range store.ImageURLs {

			// to resume an ongoing but failed in between job
			// skips images already processed
			if i < len(store.ImageUUIDs) {
				continue
			}

			img_holder, err := files.DownloadImage(img_url)

			if err != nil {
				svs.UpdateStoresVisitStatus(id, "failed", err.Error(), store.StoreID)
				return
			}

			err = img_holder.SaveImage(id.Hex(), store.StoreID)

			if err != nil {
				fmt.Println(err)
				return
			}

			// calculate perimeter
			perim := int64(img_holder.Width) * int64(img_holder.Height)

			// gpu processing simulation
			ms := 100 + rand.IntN(301)
			time.Sleep(time.Duration(ms) * time.Millisecond)

			new_image_perims[i] = perim
			new_image_uuids[i] = img_holder.ID
		}

		// store the new image perims and uuids in db
		err = svs.UpdateVisitInfo(id, visitIndex, new_image_perims, new_image_uuids)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	svs.UpdateStoresVisitStatus(id, "completed", "", "")
}
