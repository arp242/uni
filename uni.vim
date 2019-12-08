command! -range UnicodeName
			\  let s:save = @a
			\| if <count> is# -1
			\|   let @a = strcharpart(strpart(getline('.'), col('.') - 1), 0, 1)
			\| else
			\|   exe 'normal! gv"ay'
			\| endif
			\| echo system('uni -q i', @a)[:-2]
			\| let @a = s:save
			\| unlet s:save

" Simpler version which works on the current character only:
" command! UnicodeName echo
"         \ system('uni -q i', [strcharpart(strpart(getline('.'), col('.') - 1), 0, 1)])[:-2]
