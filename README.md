This repo reflects my path into understanding OpenShift's image streams.

## Repository Contents

* `app-sample` is a Go application that serves a "Hello, World!" under the `/`
endpoint.
* `appsample-manifest.yaml` contains the Deployment, Service, and Route needed
to deploy `app-sample`
* `imagestreamimport.yaml` is an `ImageStreamImport` that will create an image
stream for the `app-sample` container image (the image needs to be built and
pushed for it to work, more details below)

## How-to

This section will go through the steps involved in importing an image stream
(without copying it locally to OpenShift), and triggering a new deployment
when the source tag changes.

We'll deploy an example application. The source code is included in this
repository, under `app-sample`.

### Requirements

You can probably run the below steps in the OpenShift sandbox, but I'll be
using a cluster created on AWS with `openshift-install`.

We will use an already built version of the `app-sample`:
`quay.io/fmissi/sample`. You won't have access to write in this repository, and
we'll need to do that eventually to trigger a deployment, so you might want to
build your own version of the `app-sample` and push it to a registry namespace
you can access, then change `imagestreamimport.yaml` to reference your version
instead of mine.

We'll be using our own project, which you can create with:

```
oc new-project myapp
```
If you're on OpenShift sandbox you won't be able to do this. In this case,
you'll need to change the namespace references in the yaml files in this
repository to use the namespace you've been provided with.

### Importing the image

Before we can reference an image stream from our app Deployment we need to
create it:

```
oc create -f imagestreamimport.yaml
```

this will create the `app` image stream:

```
oc get is
NAME   IMAGE REPOSITORY                                             TAGS     UPDATED
app    image-registry.openshift-image-registry.svc:5000/myapp/app   latest   3 hours ago

oc get istag
```

We set `spec.images[0].referencePolicy.type` to `Source` so that OpenShift's
own image registry doesn't end up pulling this locally. Instead, the image is
used directly from its source (`quay.io/fmissi/sample` in this case). More on
that later.

### Enabling deployments to reference image streams

It's not possible to reference deployment from image streams out of the box.
We'll enable this directly in the image stream, but it could be done to the
deployment as well.

```
oc set image-lookup app
oc set image-lookup imagestream --list
```

This will change the image stream yaml for us, enabling it to be referenced
by other objects.

### Deploy the sample app

```
oc apply -f appsample-manifests.yaml
```

You can see it working by making a request to it:

```
curl -Lk "$(oc get route app-sample -ojson|jq -r '.spec.host')"
Hello, World!
```

Let's have a look at what image the app pods are using:

```
oc get pod -l=app=app-sample -ojson|jq '.items[0].spec.containers[0].image'
"quay.io/fmissi/sample@sha256:05d9c9067053d05ad54a8b383238b59be2c5192fb0ea49708c982e9133b98977"
```

This manifest digest belongs to the `v0.1.0` tag of `quay.io/fmissi/sample`.
The `latest` tag currently points to that same manifest.

### Enabling the deployment trigger

We want OpenShift to trigger a deployment whenever the manifest under the
`latest` tag changes, so that when a new version of `quay.io/fmissi/sample` is
imported into the image stream, OpenShift can immediately deploy it for us.

```
oc set triggers deploy/app-sample --from-image=app:latest -c app
```

### Trigering a deployment

For OpenShift to trigger a deployment two things need to happen:
1. the source image (`quay.io/fmissi/sample:latest`) needs to change
2. the image stream needs to import the new source

#### 1. Changing the source image

This is where you'll need access to the repository where the sample app image
is stored. You'll need to adapt the commands below so they work for you.

```
podman pull quay.io/fmissi/sample:v0.2.0 # you can keep this line intact
podman tag quay.io/fmissi/sample:v0.2.0 quay.io/fmissi/sample:latest
podman push quay.io/fmissi/sample:latest
```

Version `v0.2.0` of the sample app prints a "Hello, World! 123" instead of
"Hello, World!".

#### 2. Importing the new source image to the image stream

This step is the very same one we started with:

```
oc create -f imagestreamimport.yaml
```

This will update the image stream with the new reference, and that will in turn
trigger a deployment to our app-sample.

```
curl -Lk "$(oc get route app-sample -ojson|jq -r '.spec.host')"
Hello, World! 123
```

```
oc get pod -l=app=app-sample -ojson|jq '.items[0].spec.containers[0].image'
"quay.io/fmissi/sample@sha256:62b5cf5989265441a1507dc4e10f2bdb982f0ba2707487120612f055b7f2bdb8"
```

### Conclusion

We have:

* created an image stream via `ImageStreamImport`
* enabled it to be used by kubernetes objects (such as Deployment)
* deployed a sample app
* setup a deployment trigger for our app
* seen the deployment trigger at work

Hope this helps someone else wrap their head around OpenShift image streams and
deployment triggers. Writing it definitely helped me.
