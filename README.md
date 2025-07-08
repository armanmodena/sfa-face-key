Cara menjalankan
CGO_LDFLAGS="-L/usr/local/lib -ldlib -lblas -lcblas -llapack -ljpeg" CGO_CXXFLAGS="--std=c++14" go run main.go

sudo mount -t nfs -o nolock -o vers=4 192.168.3.86:`/home/webadmin/sourcode/media/sfa_mobile/face_key` /home/arman/app/sfa-face-key/faces/images