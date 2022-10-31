package graph

// type TypeChangedArgs struct {
// 	Name string
// }

// func (g *Graph) TypeChanged(ctx context.Context, args TypeChangedArgs) chan graphql.ID {
// 	rch := make(chan graphql.ID)

// 	go g.typeChanged(ctx, args.Name, rch)

// 	return rch
// }

// func (g *Graph) typeChanged(ctx context.Context, name string, rch chan graphql.ID) {
// 	defer close(rch)

// 	t := time.NewTicker(time.Second)
// 	defer t.Stop()

// 	query := map[string][]string{"domain": {name}}
// 	since := time.Now().UTC()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-t.C:
// 			bs, err := g.store.BlobNew(query, since)
// 			if err != nil {
// 				log.Printf("err: %+v\n", err)
// 			}

// 			for _, b := range bs {
// 				if b.Created.After(since) {
// 					rch <- graphql.ID(b.Hash)
// 				}
// 			}

// 			for _, b := range bs {
// 				if b.Created.After(since) {
// 					since = b.Created
// 				}
// 			}
// 		}
// 	}
// }
