package abstractions

// ResponseHandler handler to implement when a request's response should be handled a specific way.
type ResponseHandler func(response interface{}, errorMappings ErrorMappings) (interface{}, error)
