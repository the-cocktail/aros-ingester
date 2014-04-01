package reservationservice

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"github.com/fitstar/labix_mgo"
	"github.com/fitstar/labix_mgo/bson"
	"time"
)

type Reservation struct {
	Id string
	UserId string
	Total string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status string
	UserAgent string
	UserIP string
}

func log_reservation(str string, r Reservation) {
	
	log.Println(str + " : " +
	  "Id," + r.Id +
      ",UserId," + r.UserId +
  	  ",Total," + r.Total + 
	  ",CreatedAt," + r.CreatedAt.Format(time.RFC822) +
	  ",UpdatedAt," + r.UpdatedAt.Format(time.RFC822) +
	  ",Status," + r.Status +
	  ",UserAgent," + r.UserAgent +
	  ",UserIP" + r.UserIP)
}

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/reservations").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON)

	service.Route(service.GET("/{reservation-id}").To(FindReservation))
	service.Route(service.POST("").To(UpdateReservation))
	service.Route(service.PUT("/{reservation-id}").To(CreateReservation).
		Doc("Crea una reserva").
		Operation("CreateReservation").
		Reads(Reservation{}))
	service.Route(service.DELETE("/{reservation-id}").To(RemoveReservation))

	return service
}

// curl -X GET http://localhost:8080/reservations/A333DEF -H "Content-type: application/xml"
func FindReservation(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("reservation-id")
	
	session, err := mgo.Dial("localhost")
	if err != nil {	panic(err) }
	defer session.Close()
	
	c := session.DB("AROS").C("reservations")
	reserve := Reservation{}
	
	err = c.Find(bson.M{"id": id}).One(&reserve)
	if err != nil {	panic(err) }
	
	response.WriteEntity(reserve)
	log_reservation("FindReservation", reserve)
}

//curl -X PUT http://localhost:8080/reservations/AABCDEF -H "Accept: application/xml" -d"<Reservation><Id>AABCDEF</Id><UserId>2849 </UserId><Total>352</Total></Reservation>"

func CreateReservation(request *restful.Request, response *restful.Response) {
	reserve := Reservation{Id: request.PathParameter("reservation-id")}
	err := request.ReadEntity(&reserve)
	
	reserve.CreatedAt = time.Now()
	reserve.UpdatedAt = time.Now()
	reserve.Status = "E_INIT"
	
	if err == nil {
		session, err := mgo.Dial("localhost")
		if err != nil {
			panic(err)
		}
		defer session.Close()
		
		c:= session.DB("AROS").C("reservations")
		err = c.Insert(reserve)
		
		if err != nil {
			panic(err)
		}
		log_reservation("CreateReservation", reserve)
		
	} else {
		response.WriteError(http.StatusInternalServerError, err)
		log.Println(err)
	}
}

// time curl -X POST http://localhost:8080/reservations -H "Content-type: application/xml" -d"<Reservation><Id>A333DEF</Id><UserId>99999</UserId><Total>34242123123</Total></Reservation>"
func UpdateReservation(request *restful.Request, response *restful.Response) {

	reserve := Reservation{}
	err := request.ReadEntity(&reserve)
	if err != nil { panic(err) }
	
	session, err := mgo.Dial("localhost")
	if err != nil { panic(err) }

	defer session.Close()
	
	id := reserve.Id
	
	c:= session.DB("AROS").C("reservations")
	log.Println("=> Finding ID = " + id)
	
	err = c.Find(bson.M{"id": id}).One(&reserve)
	if err != nil { panic(err) }
	
	err = request.ReadEntity(&reserve)
	if err != nil { panic(err) }

	reserve.UpdatedAt = time.Now()
	
	err = c.Update(bson.M{"id": id}, reserve)
	if err != nil { panic(err) }

	response.WriteEntity(reserve)
	log_reservation("UpdateReservation", reserve)
}

func RemoveReservation(request *restful.Request, response *restful.Response) {
	
	id := request.PathParameter("reservation-id")
	reserve := Reservation{}
		
	session, err := mgo.Dial("localhost")
	if err != nil {	panic(err) }
	defer session.Close()
	
	c:= session.DB("AROS").C("reservations")
	
	err = c.Find(bson.M{"id": id}).One(&reserve)
	if err != nil {	panic(err) }
	
	err = c.Remove(bson.M{"id":id})
	if err != nil {	panic(err) }

	response.WriteEntity(reserve)
	log_reservation("RemoveReservation", reserve)
}
