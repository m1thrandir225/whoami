package handlers

type UriID struct {
	ID string `uri:"id" binding:"required"`
}
