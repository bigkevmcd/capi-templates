# capi-templates

This is a PoC for adding metadata to template parameters in CAPI objects

## Parameter metadata

We have a use for metadata for the parameters that are exposed by CAPI templates, the [SimpleProcessor](sigs.k8s.io/cluster-api/cmd/clusterctl/client/yamlprocessor) is capable of both replacing named parameters, of the form `${MY_VALUE}` and of getting the set of parameters that are defined within a template.

Irrespective of the templating system used, it would be useful to have some metadata to guide a user as to the purpose of a value they supply when rendering a template.

 * Description
 * Type of value (numeric, string, might also be possible to include sequence and map types)
 * Possible values (to allow validation, but also user-selection in a UX implementation)
 * Validation

This might be represented by a Go struct...

```go
type TemplateParame struct {
    Name string `json:"name"`
    Description string `json:"description,omitempty"`
    Type ParamType `json:"type,omitempty"` // ParamType is an enumeration of options
    Options []string `json:"options,omitempty"`
    Validations []Validation `json:"validations,omitempty"`
}
```

And a representation might look like...

```yaml
parameters:
  - name: CLUSTER_NAME
    description: Used to generate the name of the cluster and other related objects.
    type: string
    validations:
      - max 40 # I've left this a bit vague as I'm not sure what this should look like, ideally we'd reuse an existing validation schema
      - min 5
  - name: AWS_REGION
    description: Used to indicate where the cluster will be brought up.
    type: string
    options:
      - eu-west-1
      - eu-west-2
      - us-east-1
      - us-west-1
```

I've left any indication of a `default` value to the template itself as the
templating language allows for default values, and this gets messy because
technically the same parameter could have multiple different default values in the
document, this is arguably confusing for users, and a default could be introduced that because it would be provided, would indeed override the default template-wide.

## Getting the metadata into the templates

There is an existing template format, which is a stream of YAML documents, separated by the "document end" marker, with some examples in the CAPI repository here https://github.com/kubernetes-sigs/cluster-api-provider-aws/tree/main/templates.

How do we mark up templates with Metadata?

There are two proposals here, with resulting tradeoffs, adding an optional "Metadata" document, or formalising the template format.

### Metadata Document

This is the least invasive of the two options, would allow for retaining all existing templates, with a gradual addition of metadata if desired.

This involves the addition of a non-cluster document that provides the template metadata:

```yaml
---
apiVersion: cluster.x-k8s.io/v1alphaX
kind: TemplateMetadata
spec:
  description: This template is for EKS clusters.
  params:
    - name: CLUSTER_NAME
      description: This is used for the cluster naming.
---
apiVersion: cluster.x-k8s.io/v1alpha3
kind: Cluster
metadata:
  name: "${CLUSTER_NAME}"
spec:
```

A template metadata parser would enrich the set of parameters fetched by the processor's `GetVariables` call with data discovered in the `TemplateMetadata` document (if any).

This would retain full backwards compatibility with existing templates, without the metadata, the user experience would not be enhanced, but it could be incrementally added to templates as providers desire, organisations would retain access to existing templates.

On the downside, because this object is in the set of objects that are created in the cluster, it would need to be ignored by the reconciler (this is the default), there's scope for having it available separately to the main document, for example `cluster-template-eks.yaml ` and `cluster-template-eks.meta.yaml` and providing ways to fetch both files when rendering the UI which would keep the Metadata separate, this would allow for organisations to enrich templates that are provided elswehere, separation of Metadata has a cost in that changes need to be modified in two locations, which is usually more expensive than a single document, but there are options for making tooling smart about this.

### CAPITemplate Document

This is a more invasive option, existing templates would need to be converted to work with this approach.

I have modelled this after the [`TriggerTemplate`](https://tekton.dev/docs/triggers/triggertemplates/) model from the TektonCD/Triggers project.

```yaml
apiVersion: cluster.x-k8s.io/v1alphaX
kind: CAPITemplate
metadata:
  name: cluster-template
spec:
  description: this is test template 1
  params:
    - name: CLUSTER_NAME
      description: This is used for the cluster naming.
  resourcetemplates:
    - apiVersion: cluster.x-k8s.io/v1alpha3
      kind: Cluster
      metadata:
        name: "${CLUSTER_NAME}"
      spec:
        clusterNetwork:
          pods:
            cidrBlocks: ["192.168.0.0/16"]
```

This is a CR (but does not necessarily have to be actually stored in-cluster),
and the resourceTemplates are a set of templates that would be rendered and
stored in-cluster.

These files can be stored wherever you're storing your templates, but could also
be retained in-cluster, with appropriate RBAC which along with impersonation in
a Kube client, could be a strong way of controlling access to templates (if
desired).

This could be defined by the following CRD struct, with additional fields for
the other metadata items I defined above.

```go
// CAPITemplateSpec holds the desired state of CAPITemplate.
type CAPITemplateSpec struct {
	Description       string                 `json:"description,omitempty"`
	Params            []TemplateParam        `json:"params,omitempty"` // Described above
	ResourceTemplates []CAPIResourceTemplate `json:"resourcetemplates,omitempty"`
}

// CAPIResourceTemplate describes a resource to create.
type CAPIResourceTemplate struct {
	runtime.RawExtension `json:",inline"`
}

type CAPITemplate struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec holds the desired state of the CAPITemplate from the client
	// +optional
	Spec CAPITemplateSpec `json:"spec"`
	// +optional
	Status CAPITemplateStatus `json:"status,omitempty"`
}
```

It's fairly simple to enrich the data from the template parser's `GetVariables`
call with the data from the spec.

Rendering the template is slightly more complicated than it was before, as you
need to iterate over the `resourceTemplates` rather than the single `[]byte`
from reading the template, but this is providable as a method on the
`CAPITemplate` value.

As indicated above, these templates are **not** compatible with existing templates,
and organisations wanting to adopt would need to convert existing files, a simple tool that parsed the existing template file and generated a `CAPITemplate` with a set of extracted parameters all ready to be described would be trivial to implement.

I personally prefer this option, as I like the "one document" approach to
logical templates, and I think it's easier to build tooling around this
structure, but I do see the value of existing templates.
