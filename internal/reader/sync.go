package reader

func Distinction(dir1, dir2 ListFiles) ListFiles {
	if len(dir1.MapHash) == 0 && len(dir2.MapHash) == 0 {
		return ListFiles{Basepath: "", MapHash: make(map[uint64]File), MapName: make(map[string]File)}
	}

	set := ListFiles{Basepath: "", MapHash: make(map[uint64]File), MapName: make(map[string]File)}
	set.Basepath = dir1.Basepath
	for hashA, a := range dir1.MapHash {
		if b, ok := dir2.MapHash[hashA]; ok {
			//Hash found in A
			if MakeNameKey(dir1.Basepath, a.Name) == MakeNameKey(dir2.Basepath, b.Name) {
				//Name and place in the destination and source folder are the same. Means files completly identically
				entry.Infof("Files %s and %s identially - not sync", a.Name, b.Name)
			} else {
				//Names or place is different - need to sync them
				entry.Infof("Files %s and %s content is the same but place is different - add to sync", a.Name, b.Name)
				// entry.Printf("Files %s and %s content is the same but place is different - add to sync", MakeNameKey(dir1.Basepath, a.Name), MakeNameKey(dir2.Basepath, b.Name))
				set.MapName[MakeNameKey(dir1.Basepath, a.Name)] = a
			}
		} else {
			//hash not found in B
			if b, ok := dir2.MapName[MakeNameKey(dir1.Basepath, a.Name)]; ok {
				//Name and place in the destination and source folder are the same. Need to add to sync the last changed file
				if a.Info.ModTime().UnixNano() > b.Info.ModTime().UnixNano() {
					entry.Infof("Files %s and %s First file is last changed. %s > %s. Add to sync first file - A", MakeNameKey(dir1.Basepath, a.Name), MakeNameKey(dir2.Basepath, b.Name), a.Info.ModTime(), b.Info.ModTime())
					set.MapName[MakeNameKey(dir1.Basepath, a.Name)] = a
				} else if a.Info.ModTime().UnixNano() < b.Info.ModTime().UnixNano() {
					entry.Warningf("Files %s and %s Second file is last changed. %s < %s. Not added to sync due to sync A->B. last changed file is B", a.Name, b.Name, a.Info.ModTime(), b.Info.ModTime())
					// set.MapName[MakeNameKey(dir2.Basepath, b.Name)] = b
				} else {
					entry.Errorf("Files %s and %s Last change time is the same. %s == %s - CRITICAL ISSUE!", MakeNameKey(dir1.Basepath, a.Name), MakeNameKey(dir2.Basepath, b.Name), a.Info.ModTime(), b.Info.ModTime())
				}
			} else {
				//Names or place is different - need to sync them
				entry.Infof("File %s is new - add to sync", MakeNameKey(dir1.Basepath, a.Name))
				set.MapName[MakeNameKey(dir1.Basepath, a.Name)] = a
			}
		}
	}
	return set

}
