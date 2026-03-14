# RLAAS Java SDK

## Build

```bash
mvn clean package
```

## Example

```java
import io.rlaas.sdk.RlaasClient;
import io.rlaas.sdk.model.CheckRequest;
import io.rlaas.sdk.model.Decision;

public class Main {
    public static void main(String[] args) throws Exception {
        RlaasClient client = new RlaasClient("http://localhost:8080");

        CheckRequest req = new CheckRequest();
        req.setRequestId("req-1");
        req.setOrgId("acme");
        req.setTenantId("retail");
        req.setSignalType("http");
        req.setOperation("charge");
        req.setEndpoint("/v1/charge");
        req.setMethod("POST");
        req.setUserId("u1");

        Decision decision = client.checkLimit(req);
        System.out.println(decision.isAllowed() + " " + decision.getAction() + " " + decision.getReason());
    }
}
```
