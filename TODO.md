## Loading

- load rules policies and validate them (initially from file system, in the future from CRDs)
- ignore invalid rules, don't block.
- based on the rules generate the validating webbhook configuration and create/update it.
- organize rules based on the "apiVersion" & "kind"

## Execution
- when requests comes in, cast it into a AdmissionReviewRequest
- validate it against the rules loaded in the above steps
- if OK return 200 if not return 403

## Dynamic loading?
- Have a goroutine that watches the CRDs and dynamically loads the new rules

## Metrics
- expose metrics related to admission review for reporting purposes

## Report generation
- report generation. separate controller to verify how many existing object are not following policies
- this could/should be a separate binary

