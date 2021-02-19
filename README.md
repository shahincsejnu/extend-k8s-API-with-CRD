# extend-k8s-API-with-CRD



## Intuitions

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


# Resources (sequentially)

- [x] [Learning Kubernetes CRD in 2020](https://www.youtube.com/watch?v=QMRZhSpuh2w&feature=youtu.be)
- [x] [Learning Kubernetes CRDs](https://www.youtube.com/watch?v=qcSvXAxsvbY&feature=youtu.be)
- [x] [Designing CRD Types](https://www.youtube.com/watch?v=12rMmW_4rJ8&feature=youtu.be)
- [x] [How to write a Kubernetes Controller](https://www.youtube.com/watch?v=LLUhMM0cAJQ)
- [x] [Learn about code generators](https://www.openshift.com/blog/kubernetes-deep-dive-code-generation-customresources)
- [x] [Workqueue Example](https://github.com/kubernetes/client-go/tree/master/examples/workqueue)
- [x] [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)
- [x] Design of my CRD (schema, yaml)
