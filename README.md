# interceptor

This is a [Tekton](https://github.com/tektoncd/triggers) [EventListener interceptor](https://github.com/tektoncd/triggers/blob/master/examples/eventlisteners/eventlistener-interceptor.yaml) that can match on GitHub HTTP Hooks.

## pull_request events

Configured as an interceptor for `pull_request` events this picks up the `Pullrequest-Action` and `Pullrequest-Repo` headers which are provided by the eventlistener configuration.

e.g. 

```
  triggers:
    - name: dev-ci-build-from-pr
      interceptor:
        header:
        - name: Pullrequest-Action
          value: opened
        - name:  Pullrequest-Repo
          value: bigkevmcd/interceptor
        objectRef:
          kind: Service
          name: demo-interceptor
          apiVersion: v1
          namespace: cicd-environment
```

This sends the events from the `EventListener` to a `Service` named `demo-interceptor`.

If the hook is a `pull_request` event, and not an `opened` event, for the `bigkevmcd/interceptor` repo, this will fail with HTTP 212, otherwise it will return the body and an HTTP 200 response.


## push events
