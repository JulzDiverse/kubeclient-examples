# The `informer` API

The informer API makes getting kubernetes events for all resources easy. The idea is to prvoide notification functions, that are called whenever an events comes in.

There are three different types of Resources Event Handlers:

- *OnAdd:* Called whenever a resources is added
- *OnUpdate:* Called whenever a resources is modified. It is also called when a re-list happens. Therefor it is very useful for peridoical evaluations or synchronisztions.
- *OnDelete:* Called whenever a resources got deleted

## Example Code

In this example an Informer for Pods is used. It is crated by using the InformerFactory functions. There are different factory function. However, to limit the events to a single namespace the factory function with options needs to be used. The factory functions, as well as the Informer Options can be found [here](https://godoc.org/k8s.io/client-go/informers). The matching option function for namespaces is `WithNamespace(namespace)`. 

In the next step a resource event handler was added to the informer with a function for each type (`add`, `update`, `delete`). The event handler passes an object to the notification function, which usually is delegated further, but for this example it is type casted to print information about the resources. To cast the pod, the `meta` and the `core` packages are used. It makes it easy get the meta or pod information.

The intention of this example was to get the Pod IP. The best way to get it is to use the `update` event handler. The `add` is not optimal for this purpose as the Pod IP is empty if no IP was allocated yet. This is usually the case when a new resource is added.

The last step is to run the informer and provide an channel which is used to stop the informer from listing to events. 
