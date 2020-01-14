# skiplist
A fast and straightforward non-threadsafe skiplist. Supports "Get", "Set", and "Del" operations. Keys are uint64 and unique, with math.MaxUint64 used as a sentinel, which suits my current project. There is also the start of a finger implentation which is not fully fleshed out.

This is 10-20% faster and has a smaller memory footprint than all of the skiplists found at https://github.com/MauriceGit/skiplist-survey. The main reason is that this uses the "naive" skiplist structure, where each node has a right and down pointer. Other implementations use nodes with a slice of pointers, which is inefficient because of the memory overhead and additional indirection during horizontal traversal. A little math will show that there is less memory used with the naive implementation for P <= .5, and profiling confirms it.

I have also tested when compiled to .wasm, and it is about twice as fast as the next closest implementation.

One skiplist on the survey claims to be concurrent, but the operations just sit behind a global (non-rw) mutex. A simple rw mutex in front of this skiplist will perform much better. I have not been able to locate a good concurrent implementation.

~~~~
func main() {
  s := skiplist.NewSkipList()
  s.P = .25
  s.Set(123, "value")
  
  f := S.Finger()
  fmt.Println(f.Next())
  f.Reset()
  
  fmt.Println(s.Get(123))
  s.Del(123)
}
~~~~
  





