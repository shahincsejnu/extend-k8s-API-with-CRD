# extend-k8s-API-with-CRD

# Overall Process of extending k8s API with CRD

    CRD --> CR  --> Custom Controller

# CustomResourceDefinition

## CRD Structure (for `apiextensions.k8s.io/v1`)

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: <names.plural>.<spec.group>

scope: Namespaced
names:
  kind: <resource_object_name_camelcase>
  plural: <resource_object_plural_name>
  singular: <resource_object_singular_name>
  shortNames:
    - <short_name_of_resource_object>

spec:
  group: <group_name_as_your_choice>
  versions:
    - name: <version_name_as_your_choice>
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              
```

## Explanations of CRD structure fields

* `apiVersion`: The apiVersion key specifies which version of the Kubernetes API you’re using to create this object. To create a new CRD, we use “apiextensions.k8s.io/v1beta1” as the value.
* `kind`: The kind key specifies what kind of object you want to create. As we are about to create a CRD, we put “CustomResourceDefinition” as the value.
* `metadata`: The metadata key is used to define the data that can uniquely identify the object. In the example, in this tutorial, you define a name to identify the object, which is the combination of `names.plural` and `spec.group`.
* `scope`: this key determine the scope, that this object can function. There are two types of scope you can define: cluster and namspaced. If you want to manage all your resource under a certain namespace, and all of them will removed if you delete the namespace, you can choose namespaced. If you want your resource able to run in a cluster scope, which means it can only be instantiated once in one cluster, you can choose cluster.
* `names`: we use this section to define all the forms of the names for this object. The singular key determines the singular name in lowercase. The plural key determines the plural form in lowercase. The kind defines the new kind name in uppercase for this object in the cluster.
* `spec`: 
    * `group`: this key is used to specify the name of the group of this object.
    * `versions`: this key is used to define the available versions of this object. This section consists of list of name, serve, storage, schema (from v1, in v1beta1 in was in validation portion) 
        * `name`: We could have multiple versions supported at the same time. The name key specifies the name of the version.
        * `served`: The serve key specifies whether this version is still enabled in the cluster.
        * `storage`: The storage key specifies whether this version is saved in the cluster, since the cluster can save only one version.
        * `schema`: 

## Custom Resource Structure

```yaml
apiVersion:
kind:
metadata:
spec:
```


# Intuitions

- ```go
  type CronJob struct {
      replicas    *int
      states      []string
      x           map[string]string
      y           interface{}
  }
  
  var a CronJob
  b := a
  ```
  in `Go` during copying(b:= a) these things (*int(pointer), slice, map, interface) are shallow copied (means b refers `a`'s those things pointers)

- Each object of kubernetes has a method named `DeepCopy()`, which fully copied the things, like in `b := a` b and a's internal references are not same both will be differnet if we use DeepCopy()
- CRD can be registered with k8s cluster by `kubectl apply -f <crd_yaml>`
- To check the CRD you just created, run `kubectl describe crd <crd_name>`


# Code Generation for CustomResources

* deepcopy-gen—creates a method `func (t* T) DeepCopy() *T` for each type T
* client-gen—creates typed clientsets for CustomResource APIGroups
* informer-gen—creates informers for CustomResources which offer an event based interface to react on changes of CustomResources on the server
* lister-gen—creates listers for CustomResources which offer a read-only caching layer for GET and LIST requests

The last two are the basis for building controllers (or operators as some people call them). These four code-generator make up a powerful basis to build full-featured, production-ready controllers, using the same mechanisms and packages that the Kubernetes upstream controllers are using.

** During the code generator, never use extra new line (not more than one between two things, and for sequential command don't use any extra newline) **

## Code Generator Structure

    k8s-operator
        |-- pkg
        |    |-- apis
        |    |    |-- <your_CRD_group_name>
        |    |           |-- <your_CRD_version_name>
        |    |           |      |-- doc.go
        |    |           |      |-- register.go
        |    |           |      |-- types.go
        |    |           |      |-- <zz_generated_deepcopy_file>
        |    |           |-- register.go
        |    |-- <generated_client_folder>
        |           |-- clientset
        |           |-- informers
        |           |-- listers
        |-- hack
        |    |-- custom-boilerplate.go.txt
        |    |-- update-codegen.sh
        |-- vendor
        |-- go.mod
        |-- others

* for generating clientset, informers, listers: [followed this tuto](https://www.openshift.com/blog/kubernetes-deep-dive-code-generation-customresources) and took `doc.go`, `register.go`, `types.go`, `custom-boilerplate.go.txt` and `update-codegen.sh` from [here](https://github.com/openshift-evangelists/crd-code-generation) and modified the things according to my necessity.
* after creating `pkg`, `apis`, `<your_CRD_group_name>`, `<your_CRD_version_name>`, `doc.go`, `register.go`, `types.go`, and `hack` then from the `k8s-operator` run:
    - `hack/update-codegen.sh`
    - if it gives permission denied/similar error for this or some other files then make them executable by:
        - chmod +x <file_path>
* There are two kind of tags:
    * Global tags above package in doc.go
    * Local tags above a type that is processed
* Tags in general have the shape // +tag-name or // +tag-name=value, that is, they are written into comments. 
* Be prepared that an empty line might matter.


## How to generate clientset, informers, listers for your CRD


# Kubebuilder


## Generating CRDs

* KubeBuilder uses a tool called `controller-gen` to generate utility code and Kubernetes object YAML, like CustomResourceDefinitions.
* To do this, it makes use of special “marker comments” (comments that start with // +) to indicate additional information about fields, types, and packages.
* In the case of CRDs, these are generally pulled from your `_types.go` files.
* KubeBuilder provides a `make` target to run controller-gen and generate CRDs: `make manifests`.
* [all the marker comment needed for generating CRDs](https://book.kubebuilder.io/reference/markers/crd-validation.html)


# Controller

- Controllers are the core of Kubernetes, and of any operator.
- It’s a controller’s job to ensure that, for any given object, the actual state of the world matches the desired state in the object.
- Each controller focuses on one root Kind, but may interact with other Kinds.
- We call this process reconciling. In controller-runtime, the logic that implements the reconciling for a specific kind is called a Reconciler. A reconciler takes the name of an object, and returns whether or not we need to try again.

 


# Resources (sequentially)

- [x] [Learning Kubernetes CRD in 2020](https://www.youtube.com/watch?v=QMRZhSpuh2w&feature=youtu.be)
- [x] [Learning Kubernetes CRDs](https://www.youtube.com/watch?v=qcSvXAxsvbY&feature=youtu.be)
- [x] [Designing CRD Types](https://www.youtube.com/watch?v=12rMmW_4rJ8&feature=youtu.be)
- [x] [How to write a Kubernetes Controller](https://www.youtube.com/watch?v=LLUhMM0cAJQ)
- [x] [Learn about code generators](https://www.openshift.com/blog/kubernetes-deep-dive-code-generation-customresources)
- [x] [Workqueue Example](https://github.com/kubernetes/client-go/tree/master/examples/workqueue)
- [x] [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)
- [x] Design of my CRD (schema, yaml)
- [x] [Helm vs Kubernetes Operators](https://www.youtube.com/watch?v=bq8Cm-zbdqU)
- [x] [Extending Kubernetes 101](https://www.youtube.com/watch?v=yn04ERW0SbI)
- [x] [Writing Kubernetes Controllers for CRDs](https://www.youtube.com/watch?v=7wdUa4Ulwxg)
- [x] [To Crd, or Not to Crd, That is the Question](https://www.youtube.com/watch?v=xGafiZEX0YA)
- [x] https://developer.ibm.com/technologies/containers/tutorials/kubernetes-custom-resource-definitions/
- [x] [Generating CRDs](https://book.kubebuilder.io/reference/generating-crd.html)
- [x] [kubebuilder Introduction](https://book.kubebuilder.io/introduction.html)
- [x] [kubebuilder Quick Start](https://book.kubebuilder.io/quick-start.html)
- [x] [Kubebuilder Tutorial](https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html)
- [x] [Kubebuilder repo](https://github.com/kubernetes-sigs/kubebuilder)
