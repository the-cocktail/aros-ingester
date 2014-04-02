counter = 0

request = function()
	path = "/reservations/" .. counter
	wrk.method = "PUT"
	wrk.host = "localhost"
	wrk.port = 8080
	wrk.body = "<Reservation><Id>"..counter .. "</Id><UserId>"..counter.."</UserId><Total>"..counter.."</Total></Reservation>"
	wrk.headers["Accept"] = "application/xml"
	counter = counter + 1
	return wrk.format(nil, path)
end

