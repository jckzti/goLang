# goLang
Criador automático de formulários e campos para uso em HTML


# Dependencias:

No MacOS

```bash
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
brew install dep
dep init
dep ensure -add "go.mongodb.org/mongo-driver/mongo@~1.1.0"
```

```bash
brew install pkg-config dlib
sed -i '' 's/^Libs: .*/& -lblas -llapack/' /usr/local/lib/pkgconfig/dlib-1.pc
```

```bash
go get github.com/Kagami/go-face
```

Currently `shape_predictor_5_face_landmarks.dat`, `mmod_human_face_detector.dat` and
`dlib_face_recognition_resnet_model_v1.dat` are required. You may download them
from [dlib-models](https://github.com/davisking/dlib-models) repo:

```bash
mkdir models && cd models
wget https://github.com/davisking/dlib-models/raw/master/shape_predictor_5_face_landmarks.dat.bz2
bunzip2 shape_predictor_5_face_landmarks.dat.bz2
wget https://github.com/davisking/dlib-models/raw/master/dlib_face_recognition_resnet_model_v1.dat.bz2
bunzip2 dlib_face_recognition_resnet_model_v1.dat.bz2
wget https://github.com/davisking/dlib-models/raw/master/mmod_human_face_detector.dat.bz2
bunzip2 mmod_human_face_detector.dat.bz2
```