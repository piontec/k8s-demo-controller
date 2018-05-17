# Demo-Controller

`demo-controller` is the simplest, yet fully valid, kubernetes controller I could up come with. When I wanted to learn how to build k8s controllers, I search the net and found only some general ideas or already quite complicated examples, that were actually really doing "something" or were using Custom Resource Definitions (CRDs).

One of the best examples I found is `kubewatch` project by [Bitnami](https://engineering.bitnami.com/articles/kubewatch-an-example-of-kubernetes-custom-controller.html). This code is entirely based on `kubewatch` and I only removed all the parts I could to make this example as simple and as self-contained as I could. Kudos to @bitnami-labs!

## What does this controller do?

Not much. It watches for events related to lifetime of pods and logs messages about them.

## How is it different from `kubewatch`?
I removed:
* all the configuration options
* ability to watch different resource types
* some additional layers of event abstractions

That way, the controller is just about 200 LOC, plus about 50 of helpers. The rest is just "your code": direct handling of changes detected by the controller.

## How does it work? How controllers work?
In general, I think it's still pretty hard to get started with controllers and it's hard ot find a detailed documentation with examples. Reading the code is still the way to go. But let me give you some general idea how it works (and a bunch of URLs to check later).

### The controller pattern
In general, the controller works like this pseudo-code:
```
while True {
  values = check_real_values()
  tune_system(values, expected_state)
}
```
So, the controller gets the state of the system, compares it with the desired (programmed) state and does all the actions necessary to bring the current state to the desired state.

*K8s example:* Deployment controller is constantly notified about pods being removed and created. It checks how many pods with specific labels are running in the cluster right now. If this number is different from the number of replicas configured in the [Deployment object](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/), it creates or destroys pods to match their count to the desired number.

*Disclaimer:* by no means I'm an expert in topic. All the info bellow is according to my best knowledge, which is currently limited as I've just started and I'm learning as well.

### Controller in kubernetes
Despite the fact that the general idea of a controller is quite simple, the "correct" implementation in kubernetes is not. This is because there are traps and pitfalls waiting for you on the way. Fortunately, the community developed necessary tools, you "just" have to use them. So, let's have a brief look of how controller should be created to meet the best practices.

1. Use `Informer`/`SharedInformer` to get data you need from the API server. This will make sure that your results are cached and in general effectively fetched from the server. You won't have to cope with getting the data on your own and you will also spare some work on the API server. This is what code in [controller.go L71](pkg/controller/controller.go#L71) does.
  * Possible pitfall: 3rd argument to [`NewSharedIndexInformer`](https://github.com/kubernetes/client-go/blob/ea16f6128e4625e4a0377652c8704d7fd79a29de/tools/cache/shared_informer.go#L79) is the 'resynchronization period'. Resynchronization can be used to be sure that you have not missed (due to some kind of bug, disconnection or something) any updates about objects you're interested in. If you set the resynchronization period, the resynchronization will never be done. If you set it to a positive value, then every that count of seconds a resynchronization will be performed, which means **you will get 'Changed' event notification` for every object you subscribed for**, even if you already processed the change.
2. Don't act directly on events provided by the `Informer`, but queue them using [one of work queues](https://godoc.org/k8s.io/client-go/util/workqueue) and start a consumer for the queue. One of the recommended (and used here) queues is the [`RateLimitingQueue`](https://github.com/kubernetes/client-go/blob/ea16f6128e4625e4a0377652c8704d7fd79a29de/util/workqueue/rate_limitting_queue.go#L37), which limits the speed the events can be added to the queue.
  * Possible pitfall: don't forget to confirm to the queue that a message has been processed and can be forgotten: [controller.go L185](pkg/controller/controller.go#L185).
  * Possible pitfall: don't forget to cleanup after yourself when your controller is shutting down: [controller.go L141](pkg/controller/controller.go#L141).
3. Before starting your worker, make sure that the cache in your `Informer` has already filled and is in sync with API server: [controller.go L149](pkg/controller/controller.go#L149).
4. Now, happily fetch messages from your work queue and do whatever you need to do to make your controller functional.

### What to read / check next
* This [blog post](https://medium.com/@cloudark/kubernetes-custom-controllers-b6c7d0668fdf). Pay particular attention to the picture showing dependencies and data flow between different parts of the Informer-Workqueue setup. Analyze it twice. Or three times - it's worth it!
* This [Bitnami blog post](https://engineering.bitnami.com/articles/kubewatch-an-example-of-kubernetes-custom-controller.html) about the more complex version of the same logic called `kubewatch`.
* [Official documentation about building controllers](https://github.com/kubernetes/sample-controller).
* [This openshift blog entry](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/) about how to start with CRDs.
* Related topic: [Operator Pattern](https://github.com/operator-framework/getting-started) (application specific custom controller pattern).

## Building
You need `dep`. Get and install it here: [https://github.com/golang/dep](https://github.com/golang/dep). Then run,
```
# to fetch dependencies
dep ensure
# to build the whole thing
make
```

## Running
Make sure your `kubectl` is working. 

### Running as standalone binary
Just run `./demo-controller`. 

### Running as pod in a cluster
*  set `DOCKER_REPO` variable in [`Makefile`](Makefile) to point to a docker registry where you can store the image
*  run `make build-image` to build locally a docker image with your controller
*  run `make push-image` to push it to the registry
*  edit [`demo-controller.yaml`](demo-controller.yaml) and change `image: YOUR_URL:TAG` to pint to your image registry and the version tag you want to deploy
*  run `kubectl create -f demo-controller.yaml`

