(function() {
	// Write output to HTML element.
	window.set_output = (output) => {
		// Write stdout to terminal.
		let outputBuf = '';
		const decoder = new TextDecoder("utf-8");
		globalThis.fs.writeSync = (fd, buf) => {
			outputBuf += decoder.decode(buf);
			const nl = outputBuf.lastIndexOf("\n");
			if (nl != -1) {
				output.innerText += outputBuf.substr(0, nl + 1);
				window.scrollTo(0, document.body.scrollHeight);
				outputBuf = outputBuf.substr(nl + 1);
			}
			return buf.length;
		};
	};

	// Provide a readline-like input.
	window.readline = function(progname, output, input, cb) {
		var hist = [], hist_index = 0, reading_stdin = false;

		globalThis.fs.read = (fd, buffer, offset, length, position, callback) => {
			reading_stdin = true;
			output.innerText += 'reading from stdin...\n'
		};

		input.addEventListener('keydown', (e) => {
			//console.log(e.keyCode);

			// ^L
			if (e.ctrlKey && e.keyCode === 76) {
				e.preventDefault();
				output.innerText = '';
			}
			// ^P, arrow up
			else if ((e.ctrlKey && e.keyCode === 80) || e.keyCode === 38) {
				e.preventDefault();
				input.value = hist[hist.length - hist_index - 1] || '';
				if (hist_index < hist.length - 1)
					hist_index++;
			}
			// Arrow down; no ^N as it seems that can't be overridden: https://stackoverflow.com/q/38838302
			else if (e.keyCode === 40) {
				e.preventDefault();
				input.value = hist[hist.length - hist_index] || '';
				if (hist_index > 0)
					hist_index--;
			}
			// Enter
			else if (e.keyCode === 13) {
				e.preventDefault();

				if (progname !== '')
					output.innerText += '$ ' + progname + ' ';
				output.innerText += input.value + "\n";

				if (reading_stdin) {
					reading_stdin = false;
					var cmd = (hist[hist.length - 1] + ' ' + input.value).split(' ');
				}
				else {
					hist.push(input.value);
					var cmd = input.value.split(' ');
				}

				input.value = '';
				if (cmd.length === 0)
					return;
				if (cmd[0] !== progname)
					cmd = [progname].concat(cmd);

				cb(cmd);
			}
		});

		// Focus on load.
		input.focus();
	};
}());
