from datetime import timedelta

from airflow import DAG

from airflow.operators.bash import BashOperator
from airflow.utils.dates import days_ago
import configparser

from airflow.models import Variable

def read_config():
    config = configparser.ConfigParser()
    config.read('{}/conf/config.ini'.format(BASE_PATH))
    return config

def getDns(cfg):
    # fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", *databaseUser, *databasePassword, *databaseHost, *database)
    return '{}:{}@tcp({}:3306)'.format(cfg['sql']['dbUser'], cfg['sql']['dbPw'], cfg['sql']['dbHost'])

BASE_PATH = Variable.get("google-takeout-path")
config = read_config()

default_args = {
    'owner': 'airflow',
    'depends_on_past': False,
    'email_on_failure': False,
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
}
dag = DAG(
    'google-takeout',
    default_args=default_args,
    description='Clears Google Takeout Data',
    tags=['google'],
    start_date=days_ago(2),
)

'''
m1 = MysqlOperator(
    task_id='Create all tables', 
    mysql_conn_id=getDns(config), 
    sql='{}/sql/ddl.all.sql'.format(BASE_PATH), 
    dag=dag
)
'''

# didn't install the mysql dependency
# could I? yes. did I? no
m1 = BashOperator(
    task_id='create_all_tables',
    bash_command='mysql --host={dbHost} --user={dbUser} --password={dbPw}  -e "{base}/sql/ddl/all.sql"'
        .format(base=BASE_PATH, dbHost=config['sql']['dbHost'], dbUser=config['sql']['dbUser'], dbPw=config['sql']['dbPw']),
    dag=dag,
)


t1 = BashOperator(
    task_id='add_filename_to_csvs',
    bash_command='{base}/prep/add-filename-to-csv.sh --base-dir {export_dir}'.format(base=BASE_PATH, export_dir=config['fs']['baseDir']),
    dag=dag,
)

t2 = BashOperator(
    task_id='organize_files',
    depends_on_past=True,
    bash_command='{base}/prep/organize.sh --base-dir {base_dir} --target-dir {out_dir}'
        .format(base=BASE_PATH, base_dir=config['fs']['baseDir'], out_dir=config['fs']['targetDir']),
    dag=dag,
)

t3 = BashOperator(
    task_id='beam_csv_to_sql',
    depends_on_past=True,
    bash_command='''
    cd {base}pipelines/
    go run csv-to-sql.go utils.go \
	--track_csv_dir="{out_dir}/csv/*.csv" \
	--playlist_csv_dir="{out_dir}/playlists/*.csv" \
	--radio_csv_dir="{out_dir}/clean/radios/*.csv" \
	--database_host={dbHost} \
	--database_user={dbUser}\
	--database_password="{dbPw}"
    '''.format(base=BASE_PATH, out_dir=config['fs']['targetDir'], 
        dbHost=config['sql']['dbHost'], dbUser=config['sql']['dbUser'], dbPw=config['sql']['dbPw']),
    dag=dag,
)

t4 = BashOperator(
    task_id='create_flat_list_incoming',
    depends_on_past=True,
    bash_command='find {out_dir}/mp3 -type f -name "*.mp3" > all_mp3.txt'
        .format(out_dir=config['fs']['targetDir']),
    dag=dag,
)

t5 = BashOperator(
    task_id='create_flat_list_master',
    depends_on_past=True,
    bash_command='find {master_dir}/mp3 -type f -name "*.mp3" > all_mp3.txt'
        .format(master_dir=config['fs']['masterDir']),
    dag=dag,
)

t6 = BashOperator(
    task_id='beam_extract_id3v2',
    depends_on_past=True,
    bash_command='''
    cd {base}pipelines/
    go run extract_id3v2.go utils.go \
	--mp3_list="{out_dir}/all_mp3_master.txt" \
	--database_host={dbHost} \
	--database_user={dbUser} \
	--database_password="{dbPw}"
    '''.format(base=BASE_PATH, out_dir=config['fs']['targetDir'], 
        dbHost=config['sql']['dbHost'], dbUser=config['sql']['dbUser'], dbPw=config['sql']['dbPw']),
    dag=dag,
)

# Graph goes brr
m1 >> t1 >> t2 >> t3 >> t4 >> t5 >> t6