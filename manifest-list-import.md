## Importing a manifest list via ImageStreamImport

### 1. Create the ImageStreamImport

```
cat <<EOF | oc create -f -
apiVersion: image.openshift.io/v1
kind: ImageStreamImport
metadata:
  name: app
  namespace: myapp
spec:
  import: true
  images:
  - from:
      kind: DockerImage
      name: quay.io/fmissi/ubuntu
    to:
      name: latest
    referencePolicy:
      type: Source
    importPolicy:
      importMode: "ManifestList"
EOF
```

### 2. Ensure Image objects were created for the manifest list itself as well as all its sub manifests

```
oc get images|grep quay.io/fmissi/ubuntu
oc get is app
```

You may compare these against the contents of quay.io/fmissi/ubuntu:latest's manifest list.
There should be one image for the manifest list itself, and one image for each of the
sub-manifests, totalling 7 images.

### 3. Cleaning up

Delete the images previously created:
```
oc delete images $(oc get images | grep quay.io/fmissi/ubuntu | awk '{ print $1 }')
```

Delete the ImageStream (named after the ImageStreamImport):
```
oc delete is app
```
