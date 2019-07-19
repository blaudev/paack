package basic

import (
        "encoding/json"
        "errors"
        "net/http"
        "net/rpc"
        "strings"
)

const (
        integratorProcessFuncName = "Service.Process"
        maxNumOfAttempts          = 10
        apiURL                    = "http://localhost:5010/api"
        host                      = ":5002"

)

// Request represents the request to the RPC server
type Request struct {
        Customer Customer
}

// Response represents the response from the RPC server
type Response struct {
        Ok bool
}

// Service is the CMS integrator service
type Service struct {
        cli *http.Client
}

// CmsService runs a new CMS integrator service
func CmsService() error {
        sv := &Service{
                cli: &http.Client{},
        }

        if err := rpc.Register(sv); err != nil {
                return err
        }

        rpc.HandleHTTP()

        return http.ListenAndServe(host, nil)
}

// Process sends a customer to CMS API
func (sv *Service) Process(req *Request, resp *Response) error {
        data, err := json.Marshal(req.Customer)
        if err != nil {
                return err
        }

        for i := 0; i < maxNumOfAttempts; i++ {
                req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(string(data)))
                if err != nil {
                        continue
                }

                r, err := sv.cli.Do(req)
                if err != nil {
                        continue
                }

                if r.StatusCode != http.StatusOK {
                        continue
                }

                // returns if all is ok
                resp.Ok = true
                return nil
        }

        return errors.New("maximum number of attempts exceeded")
}
