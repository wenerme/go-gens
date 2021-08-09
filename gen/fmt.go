package gen

type Formatter func(file *File) error

func Formatters(f ...Formatter) Formatter {
	return func(file *File) error {
		for _, formatter := range f {
			err := formatter(file)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
