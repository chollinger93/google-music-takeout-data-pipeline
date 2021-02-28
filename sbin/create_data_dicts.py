import pandas as pd 
import glob
import argparse

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-i', '--in_dir', type=str, required=True, help='Input directory')
    args = parser.parse_args()

    df = pd.concat((pd.read_csv(f) for f in glob.glob(args.in_dir)))
    print(f'Schema for %s is:', args.in_dir)
    print(df.info(show_counts=True))
    