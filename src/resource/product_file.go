package resource

import (
	"lib"
    "resource/model"
)


func SaveProductFile(name string, categoryId int, categoryUrl string)(fileId int, err error) {
	var file = &model.ProductFile{}

	var engine = lib.Connection()
	if has,err := engine.Where("name = ?", name).And("product_type_id = ?", categoryId).Get(file); !has {
		if err != nil {
			return 0, err;
		}
	}

	file.ProductTypeId = categoryId
	file.Name = name
	file.Url = categoryUrl + "\\" + name
	file.IsDelete = 0

	filePath :=  file.Url
	covers := []string{filePath}

	getFileAttribute(covers, file)

	if file.Id > 0 {
		_, err := engine.ID(file.Id).AllCols().Update(file)
		if err != nil {
			lib.Logger().Println("update failed")
			return 0, err
		}
	} else {
		_, err := engine.InsertOne(file)
		if err != nil {
			lib.Logger().Println("insert failed")
			return 0, err
		}
	}
	
	fileId = file.Id
	return fileId, nil
}

func getFileAttribute(covers []string, file *model.ProductFile) error {
	//作品类型封面路径
	_,coverCompressPath1,coverCompressPath2 := GetCoverInfo(covers)

	file.CompressUrl1 = coverCompressPath1[0]

	file.CompressUrl2 = coverCompressPath2[0]
	return nil
}


func DeleteFile(categoryId int, activeFileIds []int) {

	var engine = lib.Connection()

	productfile := new(model.ProductFile)
	productfile.IsDelete = 1

	_, err := engine.Where("product_type_id = ?", categoryId).NotIn("id", activeFileIds).Update(productfile)
	if err != nil {
		lib.Logger().Println("Delete File Failed", err.Error())
	}
}