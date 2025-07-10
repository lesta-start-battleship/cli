package ui

func RenderError(err string) string {
	return ErrorStyle.Render("Ошибка: " + err)
}

func RenderSuccess(msg string) string {
	return SuccessStyle.Render(msg)
}
